# EniQilo Store API Documentation

## Authentication & Authorization

### Register Staff

- [ ] **POST /v1/staff/register**
  - **Input**: phoneNumber, name, password
  - **Output**: User ID, accessToken
  - **Errors**: 409 (conflict), 400 (validation), 500 (server error)

### Staff Login

- [ ] **POST /v1/staff/login**
  - **Input**: phoneNumber, password
  - **Output**: User ID, accessToken
  - **Errors**: 404 (not found), 400 (wrong password/validation), 500 (server error)

## Product Management

### Add Product

- [ ] **POST /v1/product**
  - **Input**: name, sku, category, imageUrl, notes, price, stock, location, isAvailable
  - **Output**: Product ID, createdAt
  - **Errors**: 400 (validation), 401 (token issues)

### Get Products

- [ ] **GET /v1/product**
  - **Parameters**: id, limit, offset, name, isAvailable, category, sku, price, inStock, createdAt
  - **Output**: List of products
  - **Errors**: 401 (token issues)

### Edit Product

- [ ] **PUT /v1/product/{id}**
  - **Input**: Product details
  - **Output**: Success message
  - **Errors**: 400 (validation), 401 (token issues), 404 (not found)

### Delete Product

- [ ] **DELETE /v1/product/{id}**
  - **Output**: Success message
  - **Errors**: 401 (token issues), 404 (not found)

## Customer Management

### Register Customer

- [ ] **POST /v1/customer/register**
  - **Input**: phoneNumber, name, password
  - **Output**: User ID
  - **Errors**: 400 (validation), 409 (phoneNumber exists)

### Get Customers

- [ ] **GET /v1/customer**
  - **Parameters**: phoneNumber, name
  - **Output**: List of customers
  - **Errors**: 401 (token issues)

### Checkout Products

- [ ] **POST /v1/product/checkout**
  - **Input**: customerId, productDetails, paid, change
  - **Output**: Success message
  - **Errors**: 404 (customer/product not found), 400 (validation/payment issues)

## Technical Requirements

- **Backend**: Golang, Postgres database, raw queries
- **Server**: Run on port 8080, environment variables for production
- **Docker**: Compile to Docker images, push to Docker Registry
- **Database Migration**: Use golang-migrate tool

**Note**: The above information is a concise summary of the API endpoints and requirements as described on the web page.
