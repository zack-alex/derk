**password-deriver** is a tool to derive passwords from a master password in a deterministic way. It is indended for my personal use, so breaking changes are possible and even likely.

The key derivation method is the same one that's used in the [**lesspass**](https://github.com/lesspass/lesspass).
The password encoding approach is different, though. The following procedure is used:

1. The last 12 bytes of the secret key are taken;
2. They are encoded using Base58 algorithm;
3. A single hyphen (-) is inserted at the end.
