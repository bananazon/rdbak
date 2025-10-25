package raindrop

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bananazon/raindrop/pkg/crypt"
)

func (r *Raindrop) EncryptPassword() (err error) {
	ciphertext, err := crypt.Encrypt(r.Config.Password)
	if err != nil {
		return err
	}
	r.Config.EncryptedPassword = ciphertext
	return nil
}

func (r *Raindrop) DecryptPassword() (err error) {
	plaintext, err := crypt.Decrypt(r.Config.EncryptedPassword)
	if err != nil {
		return err
	}
	r.Config.Password = plaintext
	return nil
}

// Return a list of files that is older than olderThan
func FindOldFiles(pattern string, olderThan time.Duration) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error: %w", err)
	}

	var oldFiles []string
	cutoff := time.Now().Add(-olderThan)

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("skipping %s: %v", file, err)
			continue
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			oldFiles = append(oldFiles, file)
		}
	}

	return oldFiles, nil
}

func (r *Raindrop) PruneBackupFiles(fileType string) {
	if r.PruneOlder {
		r.Logger.Infof("Looking for outdated %s backup files to prune.", fileType)

		pattern := fmt.Sprintf("%s/%s", r.RaindropRoot, fmt.Sprintf("%ss-*.yaml", fileType))
		timePeriod := 1 * Week

		oldFiles, err := FindOldFiles(pattern, timePeriod)
		if err != nil {
			r.Logger.Warnf("failed to find outdated %s backup files", fileType)
		} else {
			if len(oldFiles) > 0 {
				var filesString = "files"
				if len(oldFiles) == 1 {
					filesString = "file"
				}
				r.Logger.Infof("I found %d %s backup %s that can be pruned.", len(oldFiles), fileType, filesString)
				for _, filename := range oldFiles {
					r.Logger.Infof("Pruning %s", filename)
					err = os.Remove(filename)
					if err != nil {
						r.Logger.Warnf("Failed to delete %s", filename)
					}
				}
			} else {
				r.Logger.Infof("No outdated %s backup files found.", fileType)
			}
		}
	}
}
