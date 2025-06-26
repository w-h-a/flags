const { OpenFeature } = require("@openfeature/server-sdk");
const { OFREPProvider } = require("@openfeature/ofrep-provider");

const main = async () => {
    const provider = new OFREPProvider({
        baseUrl: "http://flags:4000",
        headers: [
            ["authorization", `Bearer mytoken`],
            ["content-type", "application/json"]
        ]
    });

    await OpenFeature.setProviderAndWait("nodejs", provider);
    
    const client = OpenFeature.getClient("nodejs");

    const contexts = [
        {},
        { targetingKey: "123456" },
        { targetingKey: "654321" },
        {}
    ];

    while (true) {
        for (const ctx of contexts) {
            const got = await client.getNumberValue("number-me", 0, ctx);

            console.log(`applying feature with ${got}`);

            await new Promise((resolve) => setTimeout(resolve, 1000));
        }
    }   
}

main();