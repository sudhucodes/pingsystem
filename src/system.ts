import os from "node:os";

function formatTime(date = new Date()) {
    return new Intl.DateTimeFormat("en-IN", {
        day: "2-digit",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: true,
    }).format(date);
}

function getSystemInfo() {
    return {
        hostname: os.hostname(),
        username: os.userInfo().username,
        platform: os.platform(),
        release: os.release(),
        time: formatTime(),
    };
}

export function startupMessage() {
    const info = getSystemInfo();

    return `
🟢 Startup

💻 Device : ${info.hostname}
👤 User   : ${info.username}
🪟 OS     : ${info.platform} ${info.release}
🕒 Time   : ${info.time}
`.trim();
}
