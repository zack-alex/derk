**derk** is a tool to derive passwords from a master password in a deterministic way.
It is indended for my personal use, so breaking changes are possible and even likely.

The key derivation method is the same one that's used in *lesspass*.
The password encoding approach is different, though.
The following procedure is used:

1. The last 12 bytes of the secret key are taken;
2. They are encoded using Base58 algorithm;
3. A single hyphen (-) is inserted at the end.

**similar projects:**
- [**lesspass**](https://github.com/lesspass/lesspass)
- [**spectre**](https://spectre.app)
- [**BIP-32 spec**](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)

**usage example:** in order to paste a password for *me@example.com* into the clibpoard,
run `echo '[{"username": "me", "domain": "example.com", "method": "v1"}]' | derk`
