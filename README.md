# This is a project for defining an arbitrary structure of documents (serialized as JSON) using JSON schema

Inspired by some CMS system.

## Validations
The schemas are validated. Additional validations must be supported. The documents are fetched using an API interface with RBAC & tag-defined access rights structure.

## Layout
There is a virtual directory tree; this is maintained by documents (type level or individual document level) parent location so the path can be queried and is cached.

## Purpose / when to use
Useful for defining a knowledgebase, website or just to store structured or semi-structured data (text, json...)

## NoSQL
The system uses a NoSQL db that stored the schema information, the layouts aswell as any additional validation rules, users and html templates. They are initialized from source code for testing purposes and initial setup.

## How to use
podman run -ti ... couchdb
go run ...

# Project layout
There is a main method that instantiates coucdb connections, upserts the relevant query definitions & example documents (incl. basic schema definitions), instantiates the routes and its handlers. The domain project

The application must serve a static website that calls the endpoints and presents the data and dynamically renders markup fragments.
