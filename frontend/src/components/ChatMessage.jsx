import React from 'react';
import ReactMarkdown from 'react-markdown';

function ChatMessage({ message }) {
    const isAi = message.role === 'ai';

    return (
        <div className={`message ${isAi ? 'ai' : 'user'}`}>
            <div className="message-content">
                <div className={`avatar ${isAi ? 'ai' : 'user'}`}>
                    {isAi ? 'A' : 'U'}
                </div>
                <div className="text-content">
                    {isAi ? (
                        <ReactMarkdown>{message.text}</ReactMarkdown>
                    ) : (
                        <div>{message.text}</div>
                    )}
                </div>
            </div>
        </div>
    );
}

export default ChatMessage;
