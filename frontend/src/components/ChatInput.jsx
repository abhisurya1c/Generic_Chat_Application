import React, { useState, useRef, useEffect } from 'react';

function ChatInput({ onSend, disabled }) {
    const [text, setText] = useState('');
    const textareaRef = useRef(null);

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!text.trim() || disabled) return;
        onSend(text);
        setText('');
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
        }
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSubmit(e);
        }
    };

    const handleChange = (e) => {
        setText(e.target.value);
        // Auto rescale height
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
            textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
        }
    };

    return (
        <div className="input-area">
            <div className="input-container">
                <textarea
                    ref={textareaRef}
                    rows={1}
                    placeholder="Send a message..."
                    value={text}
                    onChange={handleChange}
                    onKeyDown={handleKeyDown}
                    disabled={disabled}
                />
                <button
                    className="send-btn"
                    onClick={handleSubmit}
                    disabled={!text.trim() || disabled}
                >
                    <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" height="16" width="16" xmlns="http://www.w3.org/2000/svg"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
                </button>
            </div>
        </div>
    );
}

export default ChatInput;
