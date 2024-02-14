# backend code base template


## Clean Architecture
Clean Architecture is a design pattern that uses the concept of "separation of layers.
Using abstraction will make it easier for us to do unit testing, and by applying Domain-Driven Design (DDD), one applies a domain layer that holds just the business logic. 
The benefits become apparent when we refactor code to cater to changing requirements and eliminate technical debt. With Clean Architecture, we change one part of the code with minimal impact on any dependent code.
You can always mock the abstraction and change its dependencies without changing other parts of the code when it becomes obsolete or needs refactoring.
clean architecture is made of four circles

![image](assets/clean-arch-implementation.webp "CleanArchImplementation")


- Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
- Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
- Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

## Project Layers:
This project has 4 Domain layer :

    Entities (Models)
    Usecase (Repository)
    Controller (Usecase)
    Framework & Driver (Delivery)

## The diagram:

![image](assets/clean-arch.png "CleanArch")


A repository is a layer connecting the app and the unknown part outside the app. For the backend, the storage should be the part where the app talks with the database, stores the data, etc.
After that, we have a UseCase layer. It is a layer where we connect between the UI or Delivery layer and the Repository layer. This layer is used to control the data we use, and some people also call this layer a Controller layer. 
Next is the Entities layer, which stores the interface and abstraction from request, response, and data. 
Finally, we have the Delivery layer, sometimes called the handler layer. It is the layer that sends data to the outside world.



## Directory structure :

![image](assets/folder-struture.png "Folder Structure")




