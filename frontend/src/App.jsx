import React, { useState, useEffect, useRef } from 'react';
import Sidebar from './components/Sidebar';
import ChatInput from './components/ChatInput';
import ChatMessage from './components/ChatMessage';
import Login from './components/Login';
import { sendMessage, getStreamUrl, getChats, deleteChat as apiDeleteChat, getMessages } from './api';

function App() {
    const [authToken, setAuthToken] = useState(localStorage.getItem('chat_auth_token'));
    const [username, setUsername] = useState(localStorage.getItem('chat_username'));

    const [chats, setChats] = useState([]);
    const [currentChatId, setCurrentChatId] = useState(null);
    const [messages, setMessages] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const scrollRef = useRef(null);

    // Load chats from backend on auth change
    useEffect(() => {
        if (authToken) {
            loadChats();
        }
    }, [authToken]);

    const loadChats = async () => {
        try {
            const data = await getChats(authToken);
            setChats(data || []);
        } catch (err) {
            console.error("Failed to load chats", err);
            if (err.message.includes("401") || err.message.includes("Unauthorized")) {
                logout();
            }
        }
    };

    // Auto scroll
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [messages]);

    const createNewChat = () => {
        setCurrentChatId(null);
        setMessages([]);
    };

    const loadChat = async (chatId) => {
        try {
            const msgs = await getMessages(chatId, authToken);
            setCurrentChatId(chatId);
            setMessages(msgs || []);
        } catch (err) {
            console.error("Failed to load messages", err);
        }
    };

    const handleSend = async (text, useStream = true) => {
        if (!text.trim()) return;

        // Optimistic UI update
        const userMsg = { role: 'user', text, created_at: new Date().toISOString() };
        const initialMessages = [...messages, userMsg];
        setMessages(initialMessages);
        setIsLoading(true);

        try {
            if (useStream) {
                const streamUrl = getStreamUrl(text, 'llama3', currentChatId, authToken);
                const eventSource = new EventSource(streamUrl);

                const aiMsgId = Date.now();
                let aiText = "";
                let newChatId = currentChatId;

                setMessages(prev => [...prev, { role: 'ai', text: '', id: aiMsgId }]);

                eventSource.onmessage = (event) => {
                    try {
                        const data = JSON.parse(event.data);
                        // Backend sends chat_id first
                        if (data.chat_id) {
                            newChatId = data.chat_id;
                            if (!currentChatId) {
                                setCurrentChatId(newChatId);
                                loadChats(); // Refresh sidebar
                            }
                        }
                        if (data.chunk) {
                            aiText += data.chunk;
                            setMessages(prev => prev.map(msg =>
                                msg.id === aiMsgId ? { ...msg, text: aiText, role: 'ai' } : msg
                            ));
                        }
                    } catch (e) {
                        console.error("Parse error", e);
                    }
                };

                eventSource.onerror = (err) => {
                    eventSource.close();
                    setIsLoading(false);
                    if (!currentChatId && newChatId) {
                        loadChats();
                    }
                };

            } else {
                const data = await sendMessage(text, 'llama3', currentChatId, authToken);
                const aiMsg = { role: 'ai', text: data.response };
                setMessages([...initialMessages, aiMsg]);

                if (!currentChatId) {
                    setCurrentChatId(data.chat_id);
                    loadChats();
                }
                setIsLoading(false);
            }
        } catch (error) {
            console.error(error);
            setMessages(prev => [...prev, { role: 'ai', text: "Error: Failed to get response." }]);
            setIsLoading(false);
        }
    };

    const handleDeleteChat = async (chatId, e) => {
        e.stopPropagation();
        try {
            await apiDeleteChat(chatId, authToken);
            setChats(chats.filter(c => c.id !== chatId));
            if (currentChatId === chatId) {
                createNewChat();
            }
        } catch (err) {
            console.error("Failed to delete chat", err);
        }
    };

    const handleLoginSuccess = (token, user) => {
        setAuthToken(token);
        setUsername(user);
    };

    const logout = () => {
        localStorage.removeItem('chat_auth_token');
        localStorage.removeItem('chat_username');
        setAuthToken(null);
        setUsername(null);
        setChats([]);
        setMessages([]);
        setCurrentChatId(null);
    };

    if (!authToken) {
        return <Login onLoginSuccess={handleLoginSuccess} />;
    }

    return (
        <div className="app-container">
            <Sidebar
                chats={chats}
                activeChatId={currentChatId}
                onNewChat={createNewChat}
                onSelectChat={loadChat}
                onDeleteChat={handleDeleteChat}
            />
            <div className="chat-area">
                <div style={{ position: 'absolute', top: 10, right: 20, zIndex: 100 }}>
                    <small style={{ color: '#888', marginRight: 10 }}>User: {username}</small>
                    <button onClick={logout} style={{ background: 'none', border: 'none', color: '#ff4d4d', cursor: 'pointer', fontSize: '14px' }}>Logout</button>
                </div>
                <div className="messages-container" ref={scrollRef}>
                    {messages.length === 0 && !isLoading && (
                        <div style={{ display: 'flex', height: '100%', alignItems: 'center', justifyContent: 'center', color: '#888', flexDirection: 'column' }}>
                            <h2>AI Assistant</h2>
                            <p>How can I help you today?</p>
                        </div>
                    )}
                    {messages.map((msg, idx) => (
                        <ChatMessage key={idx} message={msg} />
                    ))}
                    {isLoading && !messages[messages.length - 1]?.role === 'ai' && (
                        <div className="message ai"><div className="message-content">Thinking...</div></div>
                    )}
                </div>
                <ChatInput onSend={handleSend} disabled={isLoading} />
            </div>
        </div>
    );
}

export default App;
