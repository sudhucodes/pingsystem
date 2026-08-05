import { startupMessage } from "./system.js";
import { send } from "./telegram.js";

export const events = {
    startup() {
        const message = startupMessage();
        return send(message);
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
