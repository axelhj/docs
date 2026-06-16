# API-features

## Schema level
Saving a schema - the schema is parsed to check whether it is valid
It is saved - it should have a name. The name is a content type.

## Document level
Saving a document - the document follows a schema. It is validated according to the schema and saved.
One test must ensure that the errors are comprehensible and easy.
Loading a document - simply displays it. It must be validated since the schema could have changed.

## Layout level
Render content using the layout (it will parse the associated schema and render each field according to type)
If editing, will render content using the form feature (it will parse the associated schema and render each field according to type in the form).
