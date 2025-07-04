const { OpenFeature } = require("@openfeature/server-sdk");
const { OFREPProvider } = require("@openfeature/ofrep-provider");

const tok = "mytoken"

describe("bool", () => {
    const tests = [
        {
            name: "resolve a boolean flag with TARGETING_MATCH reason",
            args: {
                apiKey: tok,
                flag: "bool_targeting_match",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "bool_targeting_match",
                value: true,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve with DEFAULT reson",
            args: {
                apiKey: tok,
                flag: "default_bool",
                defaultValue: true,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "default_bool",
                value: false,
                variant: "default",
                reason: "DEFAULT",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled",
            args: {
                apiKey: tok,
                flag: "disabled_bool",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_bool",
                value: false,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled 2",
            args: {
                apiKey: tok,
                flag: "disabled_bool_2",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_bool_2",
                value: false,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "error if we expect a boolean but get another type",
            args: {
                apiKey: tok,
                flag: "string_targeting_match",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: false,
                reason: "ERROR",
                errorCode: "TYPE_MISMATCH",
                errorMessage: "Flag is not of expected type",
                flagMetadata: {}                
            }
        },
        {
            name: "error if there is no flag",
            args: {
                apiKey: tok,
                flag: "no_such_flag",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "no_such_flag",
                value: false,
                reason: "ERROR",
                errorCode: "FLAG_NOT_FOUND",
                errorMessage: "FLAG_NOT_FOUND",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with invalid api key",
            args: {
                apiKey: "notthedroid",
                flag: "bool_targeting_match",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "bool_targeting_match",
                value: false,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with no api key",
            args: {
                apiKey: "",
                flag: "bool_targeting_match",
                defaultValue: false,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "bool_targeting_match",
                value: false,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to true when targetingKey does match",
            args: {
                apiKey: tok,
                flag: "bool_query",
                defaultValue: false,
                evalCtx: { targetingKey: "123456" }
            },
            want: {
                flagKey: "bool_query",
                value: true,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to false when targetingKey does NOT match but there is a backup rule",
            args: {
                apiKey: tok,
                flag: "bool_query",
                defaultValue: false,
                evalCtx: { targetingKey: "654321" }
            },
            want: {
                flagKey: "bool_query",
                value: false,
                variant: "false",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        }
    ];

    for (const t of tests) {
        test(t.name, async () => {
            const provider = new OFREPProvider({
                baseUrl: "http://localhost:4000",
                headers: [
                    ["authorization", `Bearer ${t.args.apiKey}`],
                    ["content-type", "application/json"]
                ]
            });

            await OpenFeature.setProviderAndWait("test", provider);
            
            const client = OpenFeature.getClient("test");

            const got = await client.getBooleanDetails(t.args.flag, t.args.defaultValue, t.args.evalCtx);
            
            expect(got).toEqual(t.want);
        });
    }
});

describe("float", () => {
    const tests = [
        {
            name: "resolve a float flag with TARGETING_MATCH reason",
            args: {
                apiKey: tok,
                flag: "float_targeting_match",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "float_targeting_match",
                value: 101.25,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve with DEFAULT reson",
            args: {
                apiKey: tok,
                flag: "default_float",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "default_float",
                value: 100.25,
                variant: "default",
                reason: "DEFAULT",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled",
            args: {
                apiKey: tok,
                flag: "disabled_float",
                defaultValue: 100.25,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_float",
                value: 100.25,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled 2",
            args: {
                apiKey: tok,
                flag: "disabled_float_2",
                defaultValue: 100.25,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_float_2",
                value: 100.25,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "error if we expect a float but get another type",
            args: {
                apiKey: tok,
                flag: "string_targeting_match",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: 0.0,
                reason: "ERROR",
                errorCode: "TYPE_MISMATCH",
                errorMessage: "Flag is not of expected type",
                flagMetadata: {}                
            }
        },
        {
            name: "error if there is no flag",
            args: {
                apiKey: tok,
                flag: "no_such_flag",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "no_such_flag",
                value: 0.0,
                reason: "ERROR",
                errorCode: "FLAG_NOT_FOUND",
                errorMessage: "FLAG_NOT_FOUND",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with invalid api key",
            args: {
                apiKey: "notthedroid",
                flag: "float_targeting_match",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "float_targeting_match",
                value: 0.0,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with no api key",
            args: {
                apiKey: "",
                flag: "float_targeting_match",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "float_targeting_match",
                value: 0.0,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 100.10 when targetingKey does match",
            args: {
                apiKey: tok,
                flag: "float_query",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "123456" }
            },
            want: {
                flagKey: "float_query",
                value: 100.1,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 100.00 when targetingKey does NOT match but there is a backup rule",
            args: {
                apiKey: tok,
                flag: "float_query",
                defaultValue: 0.0,
                evalCtx: { targetingKey: "654321" }
            },
            want: {
                flagKey: "float_query",
                value: 100.00,
                variant: "false",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
    ];

    for (const t of tests) {
        test(t.name, async () => {
            const provider = new OFREPProvider({
                baseUrl: "http://localhost:4000",
                headers: [
                    ["authorization", `Bearer ${t.args.apiKey}`],
                    ["content-type", "application/json"]
                ]
            });

            await OpenFeature.setProviderAndWait("test", provider);
            
            const client = OpenFeature.getClient("test");

            const got = await client.getNumberDetails(t.args.flag, t.args.defaultValue, t.args.evalCtx);
            
            expect(got).toEqual(t.want);
        });
    }
});

describe("integer", () => {
    const tests = [
        {
            name: "resolve an integer flag with TARGETING_MATCH reason",
            args: {
                apiKey: tok,
                flag: "int_targeting_match",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "int_targeting_match",
                value: 101,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve with DEFAULT reson",
            args: {
                apiKey: tok,
                flag: "default_int",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "default_int",
                value: 100,
                variant: "default",
                reason: "DEFAULT",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled",
            args: {
                apiKey: tok,
                flag: "disabled_int",
                defaultValue: 100,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_int",
                value: 100,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled 2",
            args: {
                apiKey: tok,
                flag: "disabled_int_2",
                defaultValue: 100,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_int_2",
                value: 100,
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "error if we expect a integer but get another type",
            args: {
                apiKey: tok,
                flag: "string_targeting_match",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: 0,
                reason: "ERROR",
                errorCode: "TYPE_MISMATCH",
                errorMessage: "Flag is not of expected type",
                flagMetadata: {}                
            }
        },
        {
            name: "error if there is no flag",
            args: {
                apiKey: tok,
                flag: "no_such_flag",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "no_such_flag",
                value: 0,
                reason: "ERROR",
                errorCode: "FLAG_NOT_FOUND",
                errorMessage: "FLAG_NOT_FOUND",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with invalid api key",
            args: {
                apiKey: "notthedroid",
                flag: "int_targeting_match",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "int_targeting_match",
                value: 0,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with no api key",
            args: {
                apiKey: "",
                flag: "int_targeting_match",
                defaultValue: 0,
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "int_targeting_match",
                value: 0,
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 101 when targetingKey does match",
            args: {
                apiKey: tok,
                flag: "int_query",
                defaultValue: 0,
                evalCtx: { targetingKey: "123456" }
            },
            want: {
                flagKey: "int_query",
                value: 101,
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 100 when targetingKey does NOT match but there is a backup rule",
            args: {
                apiKey: tok,
                flag: "int_query",
                defaultValue: 0,
                evalCtx: { targetingKey: "654321" }
            },
            want: {
                flagKey: "int_query",
                value: 100,
                variant: "false",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
    ];

    for (const t of tests) {
        test(t.name, async () => {
            const provider = new OFREPProvider({
                baseUrl: "http://localhost:4000",
                headers: [
                    ["authorization", `Bearer ${t.args.apiKey}`],
                    ["content-type", "application/json"]
                ]
            });

            await OpenFeature.setProviderAndWait("test", provider);
            
            const client = OpenFeature.getClient("test");

            const got = await client.getNumberDetails(t.args.flag, t.args.defaultValue, t.args.evalCtx);
            
            expect(got).toEqual(t.want);
        });
    }
});

describe("string", () => {
    const tests = [
        {
            name: "resolve a string flag with TARGETING_MATCH reason",
            args: {
                apiKey: tok,
                flag: "string_targeting_match",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: "fdsa",
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve with DEFAULT reson",
            args: {
                apiKey: tok,
                flag: "default_string",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "default_string",
                value: "asdf",
                variant: "default",
                reason: "DEFAULT",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled",
            args: {
                apiKey: tok,
                flag: "disabled_string",
                defaultValue: "asdf",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_string",
                value: "asdf",
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "use default if the flag is disabled 2",
            args: {
                apiKey: tok,
                flag: "disabled_string_2",
                defaultValue: "asdf",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "disabled_string_2",
                value: "asdf",
                variant: "default",
                reason: "DISABLED",
                flagMetadata: {}
            }
        },
        {
            name: "error if we expect a string but get another type",
            args: {
                apiKey: tok,
                flag: "int_targeting_match",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "int_targeting_match",
                value: "",
                reason: "ERROR",
                errorCode: "TYPE_MISMATCH",
                errorMessage: "Flag is not of expected type",
                flagMetadata: {}                
            }
        },
        {
            name: "error if there is no flag",
            args: {
                apiKey: tok,
                flag: "no_such_flag",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "no_such_flag",
                value: "",
                reason: "ERROR",
                errorCode: "FLAG_NOT_FOUND",
                errorMessage: "FLAG_NOT_FOUND",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with invalid api key",
            args: {
                apiKey: "notthedroid",
                flag: "string_targeting_match",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: "",
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to default with no api key",
            args: {
                apiKey: "",
                flag: "string_targeting_match",
                defaultValue: "",
                evalCtx: { targetingKey: "" }
            },
            want: {
                flagKey: "string_targeting_match",
                value: "",
                reason: "ERROR",
                errorCode: "GENERAL",
                errorMessage: "",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 'fdsa' when targetingKey does match",
            args: {
                apiKey: tok,
                flag: "string_query",
                defaultValue: "",
                evalCtx: { targetingKey: "123456" }
            },
            want: {
                flagKey: "string_query",
                value: "fdsa",
                variant: "true",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        },
        {
            name: "resolve to 'asdf' when targetingKey does NOT match but there is a backup rule",
            args: {
                apiKey: tok,
                flag: "string_query",
                defaultValue: "",
                evalCtx: { targetingKey: "654321" }
            },
            want: {
                flagKey: "string_query",
                value: "asdf",
                variant: "false",
                reason: "TARGETING_MATCH",
                flagMetadata: {}
            }
        }
    ];

    for (const t of tests) {
        test(t.name, async () => {
            const provider = new OFREPProvider({
                baseUrl: "http://localhost:4000",
                headers: [
                    ["authorization", `Bearer ${t.args.apiKey}`],
                    ["content-type", "application/json"]
                ]
            });

            await OpenFeature.setProviderAndWait("test", provider);
            
            const client = OpenFeature.getClient("test");

            const got = await client.getStringDetails(t.args.flag, t.args.defaultValue, t.args.evalCtx);
            
            expect(got).toEqual(t.want);
        });
    }
});