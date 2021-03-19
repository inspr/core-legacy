# Schemas and Types

A Channel Type defines the type of the information that can go through a Channel.  It is mandatory for a Channel to have a type, since it allows communication between dApps to be well defined. In this way, messages can be handled easily, and converted to both the sender and the receiver. A Channel Type basically has two attributes:

* `meta`, in which metadata is defined for the Channel Type, such as its name and annotations. For more information about the meta field, take a look [here]()
* `schema`, in which a type is properly declared.

Inspr dApps use Avro for serialization and deserialization of messages, so Schemas are defined as an Avro Schema. This allows the size of the messages exchanged to be significantly reduced. If you are not familiar with Avro, it is recommended that you take a look at the official Avro documentation for the definition of Schemas:

[Avro Schema documentation](https://avro.apache.org/docs/current/spec.html#schemas)


Messages can have types in many forms, and it is common to define them as a JSON object. There is online support for converting a JSON object to an Avro Schema, and you can check an example [here](https://toolslick.com/generation/metadata/avro-schema-from-json).

## Defining a Schema

Defining a type for a Channel is quite simple. To show how to define a Channel Type and a Schema for it, consider the following Channel:

```yaml
kind: channel
apiVersion: v1
meta:
  name: My_example_channel
  annotations:
    kafka.replication.factor: 1
    kafka.partition.number: 2
spec:
  type: channel_type_example
```
Notice that the Channel above has the type `channel_type_example`. If, for example, we want messages that pass through this Channel to be an `integer`, we can define `channel_type_example` as:

```yaml
kind: channeltype
apiVersion: v1
meta:
  name: channel_type_example
schema: {"type":"int"}
```
It is also possible to define the schema as a path to an `.avsc` or `.schema` file (which describes an Avro Schema). In this way, the defined schema will be the one written in the file. For example, if we define the schema file `my_type_schema.avsc` as:

```avsc
{"type":"int"}
```
we can redefine `channel_type_example` as:
```yaml
kind: channeltype
apiVersion: v1
meta:
  name: channel_type_example
schema: <path>/my_type_schema.avsc
```
This approach can be useful when you have two or more Channel Types that share the same Schema.
