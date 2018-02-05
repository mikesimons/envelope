# Sekrits

Sekrits is a simple envelope encryption tool designed to help any project keep their secrets in version control cheaply and securely.

The premise is very simple; we use AWS KMS to encrypt keys that we store in a keyring file.
Since everything in the keyring file is encrypted it is safe to commit to version control.

The keys in the keyring can be used to encrypt / decrypt your secrets using the sekrits tool and these can be kept next to the keyring.
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
