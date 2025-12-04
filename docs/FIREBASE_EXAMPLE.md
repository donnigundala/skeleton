# Firebase Integration Example

This example demonstrates how to use Firebase in skeleton-v2.

## Configuration

1. **Add your Firebase credentials** to `config/firebase.yaml`:

```yaml
firebase:
  # Option 1: Use a service account file
  credentials_file: "path/to/service-account.json"
  
  # Option 2: Use environment variable (recommended for production)
  # credentials_json: "${FIREBASE_CREDENTIALS}"
```

2. **Set environment variable** (if using Option 2):

```bash
export FIREBASE_CREDENTIALS='{"type": "service_account", "project_id": "...", ...}'
```

## Usage Examples

### Accessing Firebase Client

```go
package controllers

import (
    "context"
    firebase "github.com/donnigundala/dg-framework/dg-firebase"
)

type NotificationController struct {
    firebase *firebase.Client
}

func NewNotificationController(app foundation.Application) *NotificationController {
    fbClient, _ := app.Make("firebase")
    return &NotificationController{
        firebase: fbClient.(*firebase.Client),
    }
}
```

### Sending FCM Notifications

```go
import "github.com/donnigundala/dg-framework/dg-firebase/fcm"

func (c *NotificationController) SendNotification(ctx context.Context) error {
    // Get FCM client
    fcmClient, err := c.firebase.FCM(ctx)
    if err != nil {
        return err
    }
    
    // Build message using fluent API
    message := fcm.NewMessage().
        Token("device-token-here").
        Notification("Hello", "World").
        Data(map[string]string{
            "key": "value",
        }).
        Build()
    
    // Send
    response, err := fcmClient.Send(ctx, message)
    return err
}
```

### Using Firestore

```go
func (c *Controller) SaveDocument(ctx context.Context) error {
    firestoreClient, err := c.firebase.Firestore(ctx)
    if err != nil {
        return err
    }
    
    _, err = firestoreClient.Collection("users").Doc("user123").Set(ctx, map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    })
    
    return err
}
```

### Using Firebase Auth

```go
func (c *Controller) VerifyToken(ctx context.Context, idToken string) error {
    authClient, err := c.firebase.Auth(ctx)
    if err != nil {
        return err
    }
    
    token, err := authClient.VerifyIDToken(ctx, idToken)
    if err != nil {
        return err
    }
    
    // Use token.UID, token.Claims, etc.
    return nil
}
```

## Notes

- Firebase client is registered as a singleton in the DI container
- Access it via `app.Make("firebase")`
- All Firebase services (Firestore, Auth, Storage, FCM) are available through the client
- FCM has a dedicated package with fluent message builder
