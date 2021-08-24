
APP = "app"
CHANNEL = "channel"
TYPE = "type"
ALIAS = "alias"

GET = "Get"
POST = "Post"
DELETE = "Delete"
UPDATE = "Update"

QUERY_MODELS = {
    APP: {
        GET: "",
        DELETE: "",
        POST: "app",
        UPDATE: "app",
    },

    CHANNEL: {
        GET: "chname",
        DELETE: "chname",
        POST: "channel",
        UPDATE: "channel",
    },

    TYPE: {
        GET: "typename",
        DELETE: "typename",
        POST: "type",
        UPDATE: "type",
    },

    ALIAS: {
        GET: "key",
        DELETE: "key",
        POST: "",
        UPDATE: "",
    },
}