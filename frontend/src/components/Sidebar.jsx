import React from 'react';

function Sidebar({ chats, activeChatId, onNewChat, onSelectChat, onDeleteChat }) {
    return (
        <div className="sidebar">
            <button className="new-chat-btn" onClick={onNewChat}>
                <span>+</span> New Chat
            </button>

            <div className="history-list">
                {chats.map(chat => (
                    <div
                        key={chat.id}
                        className={`history-item ${chat.id === activeChatId ? 'active' : ''}`}
                        onClick={() => onSelectChat(chat.id)}
                    >
                        <span className="chat-title">{chat.title || 'New Chat'}</span>
                        <button
                            className="delete-chat-btn"
                            onClick={(e) => onDeleteChat(chat.id, e)}
                            title="Delete Chat"
                        >
                            ×
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default Sidebar;
