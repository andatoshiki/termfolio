# SSH Authentication Security

This document explains the SSH authentication modes available in termfolio and their security implications.

## Authentication Modes

Termfolio supports three authentication modes for SSH public key authentication:

### 1. `none` (Default - Recommended)

**Security Level:** ✅ Secure for public portfolios

This mode disables SSH public key authentication entirely. Users can connect via SSH without providing any keys, making it ideal for public-facing portfolio applications where anyone should be able to browse your portfolio.

```yaml
ssh:
  authMode: "none"
```

**Use when:**
- Running a public portfolio that anyone should be able to access
- You want the simplest, most open configuration
- No access restriction is needed

### 2. `authorized_keys`

**Security Level:** ✅ Secure with proper key management

This mode only allows connections from SSH public keys listed in an authorized_keys file (OpenSSH format). This is useful when you want to restrict access to specific users.

```yaml
ssh:
  authMode: "authorized_keys"
  authorizedKeys: ".ssh/authorized_keys"
```

**Use when:**
- You want to restrict access to specific trusted users
- Running an internal/private portfolio
- Need audit trail of who can connect

**Setup:**
1. Create an authorized_keys file in OpenSSH format
2. Add public keys (one per line) of authorized users
3. Set the `authorizedKeys` path in config

Example authorized_keys file:
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJw... user1@example.com
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user2@example.com
```

### 3. `allow_all`

**Security Level:** ⚠️ INSECURE - Testing Only

This mode accepts ANY SSH public key without validation. This is the legacy behavior that existed before this security fix.

```yaml
ssh:
  authMode: "allow_all"
```

**⚠️ WARNING:**
- This mode accepts connections from ANY user with ANY SSH key
- No access control or validation is performed
- Should ONLY be used for testing purposes
- NOT recommended for production use

**Use when:**
- Testing SSH connectivity
- Debugging connection issues
- **NEVER** use in production

## Configuration

### Via YAML Config File

Edit your `config.yaml`:

```yaml
ssh:
  port: 2222
  address: "0.0.0.0"
  hostKeyPath: ".ssh/host_ed25519"
  authMode: "none"  # or "authorized_keys" or "allow_all"
  # authorizedKeys: ".ssh/authorized_keys"  # Required for "authorized_keys" mode
```

### Via Environment Variables

You can override the config file using environment variables:

```bash
export SSH_AUTH_MODE="authorized_keys"
export SSH_AUTHORIZED_KEYS="/path/to/authorized_keys"
```

Available environment variables:
- `SSH_AUTH_MODE` - Override auth mode
- `SSH_AUTHORIZED_KEYS` - Override authorized keys file path
- `SSH_PORT` - Override port
- `SSH_ADDRESS` - Override listen address
- `SSH_HOST_KEY_PATH` - Override host key path

## Security Best Practices

1. **For Public Portfolios:** Use `authMode: "none"`
   - This is the most appropriate for public-facing portfolios
   - Anyone can connect without authentication
   - Simple and secure for the use case

2. **For Private/Internal Use:** Use `authMode: "authorized_keys"`
   - Restrict access to specific users
   - Maintain an authorized_keys file with trusted public keys
   - Regularly audit who has access

3. **Never Use in Production:** `authMode: "allow_all"`
   - Only for testing and debugging
   - Provides no security or access control
   - Will log security warnings

## Migration from Previous Versions

Previous versions of termfolio accepted all SSH public keys without validation (equivalent to `allow_all` mode). 

**To maintain backward compatibility** (not recommended):
```yaml
ssh:
  authMode: "allow_all"
```

**To upgrade to secure mode** (recommended):
```yaml
ssh:
  authMode: "none"
```

## Troubleshooting

### Users Can't Connect (authorized_keys mode)

1. Check that the authorized_keys file exists and is readable
2. Verify the public key format is correct (OpenSSH format)
3. Check server logs for "Loaded N authorized keys" message
4. Ensure the public key in authorized_keys matches the user's private key

### Security Warnings in Logs

If you see warnings like:
```
WARNING: SSH authentication is set to 'allow_all' mode
WARNING: This accepts ANY SSH public key without validation
```

Change your `authMode` to either `none` or `authorized_keys` for better security.

### No Keys Loaded Error

If using `authorized_keys` mode and see:
```
WARNING: Failed to load authorized keys from ...
WARNING: Rejecting all public key authentication attempts
```

1. Verify the authorized_keys file path is correct
2. Check file permissions (should be readable)
3. Ensure at least one valid key exists in the file
4. Check for syntax errors in the authorized_keys file

## Security Vulnerability Fixed

This implementation fixes a security vulnerability where the application accepted all SSH public key connections without any validation. The previous code:

```go
publicKeyAuth := func(ctx ssh.Context, key ssh.PublicKey) bool {
    return true  // ❌ Accepts any key
}
```

Has been replaced with a secure, configurable authentication system that allows administrators to control access appropriately for their use case.
