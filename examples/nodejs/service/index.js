const { OpenFeature } = require("@openfeature/server-sdk");
const { OFREPProvider } = require("@openfeature/ofrep-provider");

const main = async () => {
    const provider = new OFREPProvider({
        baseUrl: "http://flags:4000",
        headers: [
            ["authorization", `Bearer myflagstoken`],
            ["content-type", "application/json"]
        ]
    });

    await OpenFeature.setProviderAndWait("nodejs", provider);
    
    const client = OpenFeature.getClient("nodejs");

    while (true) {
        const got = await client.getBooleanValue("new-feat", false, {});

        if (got) {
            console.log(`✅ applying new feature`);
        } else {
            console.log(`❌ not applying new feature`);
        }

        await new Promise((resolve) => setTimeout(resolve, 1000));
    }
}

main();