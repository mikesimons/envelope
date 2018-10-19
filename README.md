# Envelope

Envelope is a simple envelope encryption tool designed to help any project keep their secrets in version control cheaply and securely.

## Installing
Grab an appropriate binary from the releases page or if you're on OSX `brew install mikesimons/brew/envelope`

## TLDR
In AWS go to `IAM -> Encryption keys` (bottom of left menu) and create a KMS key. Grab the full ARN of the key and use it in the place of `<KMSARN>` below. `--context` is optional but recommended to enable fine grained permissions:
```
envelope addkey --context="env=production" production awskms://<KMSARN>
```
This will create a file called `keyring.yaml` that you need to keep in version control.

```
jq '.some.secret' secrets.json | envelope encrypt --key=production
```
The encrypted secret will be printed. At the time of writing it's not possible to encrypt in-place so you'll need to replace the value of `.some.secret` in `secrets.json`.

Once you've done that you can decrypt with:
```
envelope decrypt secrets.json
```

Using the context provided when you add the key to the envelope keyring you can grant fine grained permissions using IAM. For example, the following policy will allow the given role to encrypt secrets for production but not decrypt them:
```json
{
  "Effect": "Allow",
  "Principal": {
    "AWS": "arn:aws:iam::111122223333:role/MyRole"
  },
  "Action": "kms:Encrypt",
  "Resource": "*",
  "Condition": {
    "StringEquals": {
      "kms:EncryptionContext:env": "production"
    }
  }
}
```

## How it works
The premise is very simple; we use AWS KMS to encrypt keys that we store in a keyring file.
Since everything in the keyring file is encrypted it is safe to commit to version control.

The keys in the keyring can be used to encrypt / decrypt your secrets using the envelope tool and these can be kept next to the keyring.
All you will need to decrypt is sufficient IAM permissions, the keyring file and the secrets file(s).

Features:
- Multiple data keys with encryption contexts (allowing you to fine grain permissions with IAM policies)
- Recursive decryption of structured formats (YAML, JSON & TOML)
- Blob based encryption / decryption for unstructured formats
- Asymmetric encrypt / decrypt permissions using IAM policies on KMS encryption contexts (e.g. developers could encrypt production secrets but not decrypt)
- Fine grained permissions using IAM policies on encryption contexts
- Auditing of decryption key access using AWS CloudTrail w/ KMS

If you're using EC2 then you should give your instances an instance profile capable of decrypting with they KMS keys you're using for the contexts you're encrypting secrets with.

## Similar projects
- [AWS Systems Manager Parameters](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
  - Stores and provides an API for configuration parameters with KMS support
  - Supports versioning of parameters but not storage in version control
- [SOPS](https://github.com/mozilla/sops)
  - Uses envelope encryption to encrypt YAML data
  - Supports KMS & GPG
  - Wants to encrypt entire file by default & sign that data (making it difficult to merge)
- [credstash](https://github.com/fugue/credstash) / [unicreds](https://github.com/Versent/unicreds)
  - Uses dynamodb to store data encrypted with KMS
  - Supports versioning of parameters but not storage in version control
  - No specific support for structured data so everything encrypted as blobs
