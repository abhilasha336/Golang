## Overview
basically utility function which helps to decrpt the ecrypted-ciphertext


## Index
- [Decrypt(ciphertext string,key []byte) (string, error)](#func-Decrypt)
- [Encrypt(plaintext string,key []byte) (string, error)](#func-Encrypt)
- [Hash(data string) string](#func-Hash)


### func Encrypt(plaintext string,key []byte) (string, error)
Description:

1. It takes two parameters:

- plaintext: The plaintext data that you want to encrypt as a string.
- key: The 32-byte encryption key used to secure the encryption process.

2. The function initializes an AES cipher block using the provided encryption key.

3. A new GCM (Galois/Counter Mode) instance is created using the AES cipher block. 
   GCM is a mode of operation block ciphers that provides authenticated encryption.

4. A random nonce (IV - Initialization Vector) is generated. 
   The length of the nonce is determined by gcm.NonceSize().

5. The plaintext is encrypted using the GCM seal method, which takes the nonce, 
   plaintext, and additional data (nil in this case) as input. 
   It returns the ciphertext.

6. The ciphertext is encoded in base64 format to ensure it can be safely represented as a string.

7. The function returns the base64-encoded ciphertext as a string and any error encountered 
   during the encryption  process.

### func Decrypt(ciphertext string,key []byte) (string, error)

1. It takes two parameters:

- ciphertext: The ciphertext data that you want to decrypt as a base64-encoded string.
- key: The same 32-byte encryption key that was used for encryption.

2. The function initializes an AES cipher block using the provided encryption key.

3. A new GCM (Galois/Counter Mode) instance is created using the AES cipher block. G
   CM is used to perform authenticated decryption.

4. The base64-encoded ciphertext is decoded to obtain the raw ciphertext bytes.

5. The function checks if the length of the decoded ciphertext is at least the size of the nonce 
   (determined by gcm.NonceSize()). 
6. If it's too short, an error is returned as the ciphertext is considered invalid.

7. The nonce is extracted from the beginning of the decoded ciphertext.

8. The remaining bytes of the decoded ciphertext are the actual encrypted data.

9. The GCM open method is used to decrypt the ciphertext, using the nonce, the ciphertext, and 
   additional data (nil in this case). It returns the original plaintext.

10. The function returns the decrypted plaintext as a string and any error encountered during the decryption process.


### func-Hash

   Hash(data string) string
   
Hash calculates the MD5 hash of the input data and returns it as a hexadecimal string.



