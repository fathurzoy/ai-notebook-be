package contain

const (
	ChatMessageRoleUser  = "user"
	ChatMessageRoleModel = "model"

	ChatMessageRawInitialUserPromptV1 = "You are a chatbot assistant that will answer user questions based on the provided references. You must answer using the language of the user's next chat message, even if the references are written in a different language. The references I provide will include reference numbers; however, you must never recall or mention the references using those numbers, as they are only used for raw chat sessions. This chat session is a raw session and will be reformatted later. I will provide the references before you answer. You may mention the references again if needed. You must answer with 'I don't know' if you do not have enough references."

	ChatMessageRawInitialModelPromptV1 = "Understood. I will answer your questions based solely on the provided references. I will indicate if I do not have enough information to answer. I will adapt my responses to the language you use in your subsequent turns. I will not refer to the references by their numbers.\n"
)
