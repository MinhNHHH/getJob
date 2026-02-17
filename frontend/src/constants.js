// WebRTC configuration
export const rtcConfiguration = {
    iceServers: [
        { urls: "stun:stun.l.google.com:19302" },
        { urls: "stun:stun1.l.google.com:19302" },
    ],
};

export const WEBRTC_DATA_CHANNEL = "RECORDING_CHANNEL";

// RTC message type
export const RTC_ANSWER = "Answer";
export const RTC_OFFER = "Offer";
export const RTC_CANDIDATE = "Candidate";
export const RTC_DATA_CHANNEL = "DataChannel";
export const RTC_CONNECT = "Connect";


// CONNECT message type
export const CONNECT = "Connect";
export const ACCECPT_CONNECT = "OK";