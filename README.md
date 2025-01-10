# Conference Booking System

This project is a simple conference booking system implemented in Go. It allows users to book conferences, manage waitlists, and cancel bookings.

Detailed requirement doc: https://docs.google.com/document/d/1YBVaRB4cnGRJcZYbcrSmeBc55ERlHehpoUnwMxi6D9M/edit?usp=sharing

---

## **Features**
- Add Users
- Add Conferences
- Book Conference Slots
- Manage Waitlists
- Cancel Bookings
- Automatic cleanup of expired bookings and waitlisted candidates

---

## **API Documentation**
The API endpoints are provided in the Postman collection below.

### **Postman Collection**
You can import the Postman collection into your Postman workspace.

1. Save the content below as `postman-collection.json`.

```json
{
  "info": {
    "name": "Conference Booking System",
    "description": "API collection for managing users, conferences, and bookings.",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Add User",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"id\": \"user1\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/user",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["user"]
        }
      },
      "response": []
    },
    {
      "name": "Add Conference",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"TechConf2025\",\n  \"start_time\": \"2025-01-15T10:00:00Z\",\n  \"end_time\": \"2025-01-15T20:00:00Z\",\n  \"available_slots\": 100\n}"
        },
        "url": {
          "raw": "http://localhost:8080/conference",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["conference"]
        }
      },
      "response": []
    },
    {
      "name": "Book Conference",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"conference_name\": \"TechConf2025\",\n  \"user_id\": \"user1\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/booking",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["booking"]
        }
      },
      "response": []
    },
    {
      "name": "Get Booking Status",
      "request": {
        "method": "GET",
        "url": {
          "raw": "http://localhost:8080/booking/{id}",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["booking", "{id}"],
          "variable": [
            {
              "key": "id",
              "value": ""
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Cancel Booking",
      "request": {
        "method": "DELETE",
        "url": {
          "raw": "http://localhost:8080/booking/{id}",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["booking", "{id}"],
          "variable": [
            {
              "key": "id",
              "value": ""
            }
          ]
        }
      },
      "response": []
    }
  ]
}
