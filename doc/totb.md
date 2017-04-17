
## Time-based One Time Passwords
Time-based one time passwords (RFC 6238) have gotten some traction for multi-factor authentication with the "Google Authenticator" mobile app. Reviewing the RFC for a secure implementation is highly recommended, but this might make it easier.

### Options
    LookAhead: 1           // Allow passwords from n future 30-second blocks
    LookBehind: 1          // Allow passwords from n previous 30-second blocks
    B32Blocks: 8           // Secret length (change will invalidate stored pws)
    HyphB32: true          // Hyphenate base 32 encoded secrets