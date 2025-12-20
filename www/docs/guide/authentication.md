# Authentication & Users

Manage users and authentication in your Mist instance.

## User Roles

Mist supports two user roles:

### Admin

- Full system access
- Create and manage users
- Access all projects and applications
- Configure system settings
- View audit logs

### User

- Create and manage own projects
- Deploy applications
- Invite members to projects
- Limited to assigned projects

## Authentication

### Login

Access Mist at `http://your-server:3000`:

1. Enter email and password
2. Click **"Sign In"**
3. JWT token stored in HTTP-only cookie

### Session Management

- Sessions expire after 7 days of inactivity
- Tokens stored in secure HTTP-only cookies
- Automatic refresh on activity

### First-Time Setup

On first access, create admin account:

1. Enter admin email
2. Set strong password
3. Click **"Create Admin Account"**

## User Management

### Creating Users (Admin Only)

1. Go to **Settings** → **Users**
2. Click **"New User"**
3. Enter email and password
4. Select role (Admin or User)
5. Click **"Create"**

### Updating Users

Edit user details:
- Update email
- Reset password
- Change role

### Deleting Users

Remove users from the system:

1. Navigate to user list
2. Click delete icon
3. Confirm deletion

::: warning
Deleting a user removes their access but preserves their projects and applications.
:::

## Password Security

- Minimum 8 characters
- Hashed with bcrypt
- Salted for additional security

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Authentication Enhancements</h4>
</div>

- **Password Reset** - Email-based password recovery
- **Email Verification** - Verify email addresses
- **Two-Factor Authentication** - TOTP/Authenticator app support
- **SSO Integration** - OAuth2, SAML support
- **API Tokens** - Generate tokens for CLI/API access
- **Session Management UI** - View and revoke active sessions
- **User Invitations** - Invite users via email
- **Fine-grained Permissions** - Custom permission sets
