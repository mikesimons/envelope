# Envelope [![Build Status](https://travis-ci.org/mikesimons/envelope.svg?branch=master)](https://travis-ci.org/mikesimons/envelope)
Envelope is a tool designed to help any project keep their secrets in version control cheaply and securely.

## How it works
Envelope uses AWS KMS to generate and encrypt "data keys" that are stored in a keyring file and these data keys are used to encrypt / decrypt your data.

Since everything in the keyring file is encrypted it is safe to commit to version control.

All you will need to decrypt is sufficient IAM permissions, the keyring file and the secrets file(s).

Features:
- "Profiles" which allow multiple data keys with differing KMS keys and encryption contexts
- Recursive decryption of structured formats (YAML, JSON & TOML)
- Blob based encryption / decryption for unstructured formats
- Asymmetric encrypt / decrypt permissions using IAM policies on KMS encryption contexts (e.g. developers could encrypt production secrets but not decrypt)
- Fine grained permissions using IAM policies on encryption contexts
- "Encrypt in place" functionality for ease of use
- Auditing of decryption key access using AWS CloudTrail w/ KMS

## Installing
Grab an appropriate binary from the releases page or if you're on OSX `brew install mikesimons/brew/envelope`

## First steps
In AWS go to `IAM -> Encryption keys` (bottom of left menu) and create a KMS key.
Grab the full ARN of the key and use it in the place of `<KMSARN>` below.
`--context` is optional but recommended to enable fine grained permissions:

```
envelope profile add --context="env=dev" dev awskms://<KMSARN>
```
This will create a file called `keyring.yaml` that you need to keep in version control.

Given the following configuration file in `config.yaml`:
```
myservice:
  database_username: user
  database_password: pass
````

You can run the following command to encrypt the password:
```
envelope encrypt --profile dev --key myservice.database_password --in-place config.yaml
```

This will write the encrypted value directly to the file. Check the [limitations](#limitations) section for caveats around encrypting in-place.

It is also possible to encrypt entire files as blobs with envelope:
```
cat config.yaml | envelope encrypt --profile dev > config.yaml.enc
```

Or to set previously unset keys:
```
echo "somevalue" | envelope encrypt --profile dev --key my.new.key config.yaml
```

To see the decrypted values of these files:
```
envelope decrypt config.yaml
```
or
```
envelope decrypt --format=blob config.yaml.enc
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

You must take special care around providing decryption access to users with contexts.
If you create a policy that allows users to decrypt without a context condition they will be able to decrypt *ALL* values.

## Limitations
### Numeric & boolean values will get converted to strings when processed by envelope
Due to the fact that the golang yaml parser will default all scalar values to strings when unmarshalling in to interface{} maps there is no way to avoid this right now.

### Map keys will be lexicographically sorted in output when processed by envelope
Due to the fact that YAML & JSON marshallers internally use golang maps for k/v structures and golang maps have non-deterministic ordering of keys the marshallers will sort map keys when emitting marshalled output. This means that so does envelope.

### Comments and formatting are stripped when processed by envelope
Since the structured parsers do not retain comment nor formatting information it is not currently possible to preserve these when processing files with envelope.

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
