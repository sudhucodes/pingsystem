import { events } from "./events.js";

async function main() {
    console.log("BootPing started");

    await events.startup();

    // Examples
    // await events.shutdown();
    // await events.sleep();
    // await events.wake();
    // await events.custom("Hello from BootPing 🚀");
}

main().catch(console.error);
