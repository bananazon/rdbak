package raindrop

import "github.com/gdanko/rdbak/pkg/crypt"

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
