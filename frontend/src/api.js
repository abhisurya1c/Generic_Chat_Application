const API_BASE = "http://localhost:8080/api";

export async function sendMessage(prompt, model = "llama3", chatId = 0, token) {
    try {
        const response = await fetch(`${API_BASE}/chat`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ prompt, model, chat_id: chatId }),
        });

        if (!response.ok) {
            const data = await response.json().catch(() => ({}));
            throw new Error(data.error || response.statusText);
        }

        return response.json();
    } catch (error) {
        console.error("API call failed:", error);
        throw error;
    }
}

export function getStreamUrl(prompt, model = "llama3", chatId = 0, token) {
    const encodedPrompt = encodeURIComponent(prompt);
    const encodedModel = encodeURIComponent(model);
    const encodedChatId = chatId ? chatId : 0;
    return `${API_BASE}/chat/stream?prompt=${encodedPrompt}&model=${encodedModel}&chat_id=${encodedChatId}&token=${token}`;
}

export async function getChats(token) {
    const response = await fetch(`${API_BASE}/history/chats`, {
        headers: { "Authorization": `Bearer ${token}` }
    });
    if (!response.ok) throw new Error(response.statusText);
    return response.json();
}

export async function getMessages(chatId, token) {
    const response = await fetch(`${API_BASE}/history/messages?chat_id=${chatId}`, {
        headers: { "Authorization": `Bearer ${token}` }
    });
    if (!response.ok) throw new Error(response.statusText);
    return response.json();
}

export async function deleteChat(chatId, token) {
    const response = await fetch(`${API_BASE}/history/delete?chat_id=${chatId}`, {
        method: "DELETE",
        headers: { "Authorization": `Bearer ${token}` }
    });
    if (!response.ok) throw new Error(response.statusText);
}
