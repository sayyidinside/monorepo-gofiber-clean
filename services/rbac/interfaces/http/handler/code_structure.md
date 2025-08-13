### Overview of the Code Structure

This code defines a `PermissionHandler` interface and implements it using the `permissionHandler` struct, which manages **permission-related** HTTP requests in a GoFiber-based application. It leverages the service layer and the helpers to standardize request handling, response formatting, logging, validation, and error handling.

### Key Elements

1. **Interface Definition**:
   The `PermissionHandler` interface defines the contract for the handler, ensuring the implementation provides functions to handle:

   - `GetPermission`: Fetch a specific permission by its ID.
   - `GetAllPermission`: Fetch all permissions with query options.
   - `CreatePermission`: Create a new permission.
   - `UpdatePermission`: Update an existing permission by ID.
   - `DeletePermission`: Delete a permission by ID.

   This is used to decouple the implementation from the actual routing and to make the code more modular and easier to test.

2. **Struct Implementation**:
   The `permissionHandler` struct implements the `PermissionHandler` interface and depends on a `PermissionService` (injected through the constructor), allowing it to interact with the domain logic.

   ```go
   type permissionHandler struct {
       service service.PermissionService
   }
   ```

3. **Constructor Function**:
   The `NewPermissionHandler` function initializes a new `permissionHandler` by injecting the `PermissionService` dependency, enabling dependency injection for better testing and separation of concerns.

   ```go
   func NewPermissionHandler(service service.PermissionService) PermissionHandler {
       return &permissionHandler{
           service: service,
       }
   }
   ```

4. **Handler Methods**:
   Each method corresponds to an HTTP route handler:

   - **GetPermission**: Parses the permission ID from the URL, checks for errors, and then fetches the permission using the service layer. It also handles response formatting and logging.

   - **GetAllPermission**: Handles query parsing, sanitization, and fetches all permissions with pagination support. The base URL and query parameters are used to construct pagination links.

   - **CreatePermission**: Parses the request body, validates the input, and sends the data to the service layer for creation.

   - **UpdatePermission**: Similar to `CreatePermission`, but also parses the permission ID and sends the request to update an existing permission.

   - **DeletePermission**: Handles the deletion of a permission by its ID and sends the result back to the client.

   In each method, error handling is done with proper response formatting, and logs are created using `helpers.CreateLog(c)` to track each request and response cycle.

### Notable Practices

- **Error Handling**: If errors are encountered (e.g., invalid ID format or malformed body), they are properly handled with error messages and a `400 Bad Request` status.
- **Logging**: Each request is logged by calling the `helpers.CreateLog(c)` function to create structured logs that can later be used for auditing or debugging.
- **Validation**: The `ValidateInput` function ensures that incoming data is validated before being processed further.
- **Response Formatter**: The `ResponseFormatter` function is used consistently to format responses with a standard structure, including success, message, and any additional errors or data.

### Conclusion

This code follows **clean architecture principles** by separating concerns into distinct layers:

- The **handler layer** manages HTTP requests and responses.
- The **service layer** (injected through `PermissionService`) contains the business logic.
- **Helper functions** manage logging, error handling, and response formatting, ensuring consistent behavior across the entire application.

This structure promotes modularity, testability, and maintainability in your GoFiber application.
