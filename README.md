![gocker](https://user-images.githubusercontent.com/46195831/88464597-129dfb80-ceb4-11ea-8b75-fe298ac02fb6.jpeg)
# Blog-Api
This is a dockerized API for a subscription-based blogging application written in Golang.

## How to run application
  You can run the application: <br/>
  
  -Locally, by cloning this repo to your computer using this command `git clone https://github.com/Ghvstcode/Blog-Api.git`, navigating into the project directory and running `go run main.go` <br/>
  
  -From a docker image, by running `docker run -d ghvst/blog-api` this will startup the container in the image and you can make requests to `http://localhost:8080` <br/>
  You can also build the image using the docker-compose.yml file by running `docker-compose up`
  ## Endpoints
  - Route: /logs <br/>
  Method: Get <br/>
  Function: View logs <br/>
  
  - Route: api/user/new <br/>
  Method: Post <br/>
  Req.Body: Username, Email, password <br/>
  Function: Create a new user <br/>
  
  - Route: api/user/login <br/>
  Method: Post <br/>
  Req.Body:Email, password <br/>
  Function: login a user <br/>
  
  - Route: api/user/posts <br/>
  Method: Get <br/>
  Function: Used to fetch all posts a user authored <br/>
  Auth: True <br/>
  
  - Route: api/user/resetPassword <br/>
  Method: Post <br/>
  Function: Reset a users password v
  Auth: True <br/>
  
  - Route: api/user/feed <br/>
  Method: Get <br/>
  Function: Curated feed for a user <br/>
  Auth: True <br/>
  
  - Route: api/blog/new <br/>
  Method: Post <br/>
  Function: Create a new blog post <br/>
  Auth: True <br/>
  
  - Route: api/blog/{id} <br/>
  Method: Put <br/>
  Function: Update a blog post <br/>
  Auth: True <br/>
  
  - Route: api/blog/{id} <br/>
  Method: Delete <br/>
  Req.Body: 	ID,Title,Content,Author,Published, Paid,Price <br/>
  Function: Delete a blog post blog post <br/>
  Auth: True <br/>
  
  - Route: api/blog/{id} <br/>
  Method: Get <br/>
  Function: Get a specific post that you created or is free/ a user subscribed to and is published <br/>
  Auth: True <br/>
  
  - Route: api/blog/{id}/subscribe <br/>
  Method: Post <br/>
  Function: Subscribe to a particular blog post <br/>
  Auth: True<br/>
  
  ## TO-DO
  - [ ] Fix the feed route<br/>
  - [ ] Add Payment gateway to the subscribe route
  - [ ] Write Tests
  
   ## Author
   GhvstCode
