import { send } from "./telegram.js";

export const events = {
    startup() {
        return send("🟢 Startup Event");
    },

    shutdown() {
        return send("🔴 Shutdown Event");
    },

    sleep() {
        return send("😴 Sleep Event");
    },

    wake() {
        return send("☀️ Wake Event");
    },

    custom(message: string) {
        return send(message);
    },
};
