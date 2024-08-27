## 🗨️ Forum Web Application 🗨️

Welcome to the Forum Web Application project! This document will guide you through the setup and features of the forum, including how to use Docker to run your application, and the key functionalities you'll find in this web forum.
🎯 Objectives

The Forum Web Application enables:

    Communication between users.
    Categorization of posts.
    Liking and Disliking posts and comments.
    Filtering posts based on different criteria.

💾 Database

We use SQLite to manage our data, including users, posts, comments, and more. SQLite provides an efficient way to handle local data storage. Key SQL operations include:

    CREATE: To define new tables.
    INSERT: To add new records.
    SELECT: To query and retrieve data.

Refer to the SQLite Documentation for more details.
Database Structure

Based on the entity-relationship diagram, you should structure your database for optimal performance. Consider designing tables for:

    Users: to store user information.
    Posts: to store forum posts.
    Comments: to store user comments.
    Categories: to associate posts with categories.
    Likes/Dislikes: to track user feedback on posts and comments.

🔐 Authentication
User Registration

Users can register with:

    Email: Unique identifier for each user.
        If the email is already taken, an error will be returned.
    Username: Unique identifier for display purposes.
    Password: Must be encrypted (Bonus task).

Login

Registered users can log in to create posts and comments. Sessions are managed via cookies, ensuring:

    Single Session: One active session per user.
    Expiration: Each session has a defined lifespan.
    UUID: Optional, but recommended for session management.

💬 Communication
Posts & Comments

    Creation: Only registered users can create posts and comments.
    Visibility: All users can view posts and comments.
    Categories: Users can associate posts with one or more categories.

👍👎 Likes and Dislikes

    Registered Users Only: Only users who are logged in can like or dislike posts and comments.
    Visibility: Like and dislike counts are visible to all users.

🔍 Filter

Registered users can filter posts by:

    Categories: Display posts within selected categories.
    Created Posts: Show posts created by the logged-in user.
    Liked Posts: Filter posts liked by the logged-in user.

🐳 Docker

This project uses Docker for containerization. To run your forum application with Docker:

    Install Docker: Follow the Docker Installation Guide.
    Build the Docker Image: Run docker build -t forum-app . in the project directory.
    Run the Container: Use docker run -p 80:80 forum-app to start the application.
## Usage/Examples

```go

go run main.go # this will start the web
```



## Authors

- [@emarei](https://www.github.com/emarei)
- [@yyahya](https://www.github.com/yyahya)
- [@falsayya](https://www.github.com/falsayya)
- [@oabdulra](https://www.github.com/oabdulra)
- [@mohani](https://www.github.com/mohani)

```