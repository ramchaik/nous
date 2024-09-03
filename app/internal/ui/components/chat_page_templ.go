// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "nous/internal/models"

func ChatPage(sessionID, chatID string, chats []*models.Chat, query string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"chat-messages\" class=\"flex-grow overflow-y-auto mb-4 space-y-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, chat := range chats {
				if chat.Type == "user" {
					templ_7745c5c3_Err = UserMessage(chat.Text).Render(ctx, templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = BotMessage(chat.Text).Render(ctx, templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"overlay\" class=\"hidden fixed inset-0 bg-gray-900 bg-opacity-50 z-50 flex justify-center items-center\"><div class=\"bg-gray-800 rounded-lg p-6 flex flex-col items-center\"><div class=\"animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-primary-500 mb-4\"></div><p class=\"text-primary-300 font-semibold\">Thinking...</p></div></div><div class=\"w-full max-w-2xl mx-auto\"><form id=\"chat-form\" hx-post=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs("/chat/" + chatID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/ui/components/chat_page.templ`, Line: 25, Col: 31}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#chat-messages\" hx-swap=\"beforeend\" class=\"relative\"><input type=\"hidden\" name=\"sid\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(sessionID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/ui/components/chat_page.templ`, Line: 30, Col: 53}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input id=\"chat-input\" type=\"text\" name=\"query\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(query)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/ui/components/chat_page.templ`, Line: 35, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" placeholder=\"Ask me anything...\" class=\"w-full py-4 px-6 bg-gray-800 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 transition duration-300 shadow-lg\" required> <button id=\"submit-button\" type=\"submit\" class=\"absolute right-3 top-1/2 transform -translate-y-1/2 bg-primary-600 text-white p-2 rounded-full hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 transition duration-300\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-6 w-6\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z\"></path></svg></button></form></div><script>\n            // Function to get URL parameters\n            function getUrlParameter(name) {\n                name = name.replace(/[\\[]/, \"\\\\[\").replace(/[\\]]/, \"\\\\]\");\n                var regex = new RegExp(\"[\\\\?&]\" + name + \"=([^&#]*)\");\n                var results = regex.exec(location.search);\n                return results === null\n                    ? \"\"\n                    : decodeURIComponent(results[1].replace(/\\+/g, \" \"));\n            }\n\n            // Function to remove URL parameters\n            function removeUrlParameters() {\n                window.history.replaceState(\n                    {},\n                    document.title,\n                    window.location.pathname,\n                );\n            }\n\n            // Scroll to bottom after new messages are added\n            document.body.addEventListener(\"htmx:afterSwap\", function (event) {\n                var chatMessages = document.getElementById(\"chat-messages\");\n                chatMessages.scrollTop = chatMessages.scrollHeight;\n                addCopyListeners();\n            });\n\n            // Show loader on form submission\n            document\n                .getElementById(\"chat-form\")\n                .addEventListener(\"htmx:beforeRequest\", function (event) {\n                    document.getElementById(\"chat-input\").value = \"\";\n                    document\n                        .getElementById(\"submit-button\")\n                        .classList.add(\"hidden\");\n                    document\n                        .getElementById(\"overlay\")\n                        .classList.remove(\"hidden\");\n                });\n\n            // Hide loader after response is received\n            document.body.addEventListener(\"htmx:afterSwap\", function (event) {\n                document\n                    .getElementById(\"submit-button\")\n                    .classList.remove(\"hidden\");\n                document.getElementById(\"overlay\").classList.add(\"hidden\");\n            });\n\n            // Function to remove all \"Copied!\" spans\n            function removeCopiedSpans() {\n                document\n                    .querySelectorAll(\".copied-span\")\n                    .forEach((span) => span.remove());\n            }\n\n            // Function to add copy listeners to all copy buttons\n            function addCopyListeners() {\n                document.querySelectorAll(\".copy-btn\").forEach((btn) => {\n                    btn.addEventListener(\"click\", function () {\n                        removeCopiedSpans();\n                        const text =\n                            this.parentElement.querySelector(\"p\").textContent;\n                        navigator.clipboard.writeText(text).then(() => {\n                            const copiedSpan = document.createElement(\"span\");\n                            copiedSpan.textContent = \"Copied!\";\n                            copiedSpan.className =\n                                \"copied-span text-xs absolute top-2 right-8 bg-gray-800 text-white px-2 py-1 rounded\";\n                            this.parentElement.appendChild(copiedSpan);\n                            setTimeout(() => copiedSpan.remove(), 1000);\n                        });\n                    });\n                });\n            }\n\n            addCopyListeners();\n\n            // Add event listener for Cmd+K\n            document.addEventListener(\"keydown\", function (e) {\n                if ((e.metaKey || e.ctrlKey) && e.key === \"k\") {\n                    e.preventDefault();\n                    document.getElementById(\"chat-input\").focus();\n                }\n            });\n\n            // Handle initial load, URL parameters, and form submission\n            window.addEventListener(\"load\", function () {\n                const chatForm = document.getElementById(\"chat-form\");\n                const chatInput = document.getElementById(\"chat-input\");\n                const sessionIdInput =\n                    document.querySelector('input[name=\"sid\"]');\n\n                const urlQuery = getUrlParameter(\"query\");\n                const urlSessionId = getUrlParameter(\"sid\");\n\n                const storedSessionId = sessionStorage.getItem(\n                    \"sessionId\",\n                    urlSessionId,\n                );\n                const chatExists = sessionStorage.getItem(\"chatExists\");\n                const makeRequest = !storedSessionId && !chatExists;\n                if (!makeRequest) {\n                    return;\n                }\n\n                if (urlSessionId) {\n                    sessionStorage.setItem(\"sessionId\", urlSessionId);\n                    sessionIdInput.value = urlSessionId;\n                } else {\n                    if (storedSessionId) {\n                        sessionIdInput.value = storedSessionId;\n                    }\n                }\n\n                if (urlQuery) {\n                    chatInput.value = urlQuery;\n                    removeUrlParameters();\n\n                    // Submit the form\n                    const formData = new FormData(chatForm);\n                    fetch(chatForm.getAttribute(\"hx-post\"), {\n                        method: \"POST\",\n                        body: formData,\n                    })\n                        .then((response) => response.text())\n                        .then((html) => {\n                            document\n                                .getElementById(\"chat-messages\")\n                                .insertAdjacentHTML(\"beforeend\", html);\n                            chatInput.value = \"\";\n                            document.getElementById(\"chat-messages\").scrollTop =\n                                document.getElementById(\n                                    \"chat-messages\",\n                                ).scrollHeight;\n                            addCopyListeners();\n                        })\n                        .catch((error) => console.error(\"Error:\", error))\n                        .finally(() => {\n                            sessionStorage.setItem(\"chatExists\", true);\n                            document\n                                .getElementById(\"submit-button\")\n                                .classList.remove(\"hidden\");\n                            document\n                                .getElementById(\"overlay\")\n                                .classList.add(\"hidden\");\n                        });\n                }\n            });\n        </script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = BaseLayout("Nous - Chat").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate