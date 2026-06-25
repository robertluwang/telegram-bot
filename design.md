Master Guide: The Legacy iOS AI Architecture
Reviving obsolete iPads using Telegram, Golang, Tailscale, and LiteLLM.
The Problem
Older iPads (running iOS 9 through iOS 12) have perfectly functional hardware but are artificially bricked by outdated web browsers. Because Safari's WebKit engine crashes on modern frameworks ("WebKit encountered an internal error"), running native AI web apps like ChatGPT or Gemini is impossible. Furthermore, modern App Store applications require iOS 14+.
The Solution
By utilizing a microservice architecture over a Zero Trust network, we bypass the iPad's limitations entirely.
1.	The UI: We use an older version of Telegram as the graphical frontend natively on the iPad.
2.	The Frontend: A hyper-efficient Golang script listens to Telegram and routes messages. This can run on a Tiny VM or locally on the iPad using the iSH emulator.
3.	The Backend (VM 2): A dedicated LiteLLM gateway processes the AI requests.
4.	The Network: Both the frontend and backend are locked down inside an encrypted Tailscale WireGuard mesh.
Phase 1: The Client Interface (iPad)
1. Install Telegram & iSH
1.	On your iOS 12 iPad, open the App Store.
2.	Go to the Purchased tab.
3.	Find and download Telegram and iSH Shell.
4.	(Note: You must have "purchased" these previously on a newer device using the same Apple ID).
2. Create the Bot
1.	Open Telegram and search for BotFather (the official Telegram bot creator).
2.	Send the command /newbot.
3.	Follow the prompts to name your bot and give it a username (must end in bot).
4.	BotFather will give you an API Token (e.g., 123456789:ABC-DEF1234ghIkl-zyx57W2v1u123ew11). Save this.
Phase 2: The Backend (LiteLLM on VM 2)
Your dedicated VM runs LiteLLM. To prevent unauthorized access, we will bind it strictly to the Tailscale network.
1. Get the Tailscale IP
Run the following command on VM 2:
Note your VM's Tailscale IP (it starts with 100.x.x.x).
2. Start LiteLLM Securely
Start your LiteLLM server, ensuring you bind the host to the Tailscale IP rather than 0.0.0.0 (which would expose it to the public internet).
(Tip: Run this in a tmux session or as a systemd service).
3. Generate a Virtual Key
1.	Open your LiteLLM Admin UI.
2.	Generate a new Virtual Key (e.g., sk-1234abcd...). This ensures you never expose your raw Google/OpenAI keys to your other servers.
Phase 3: The Frontend Bridge (Golang)
We use Golang because it compiles to a single, zero-dependency binary and runs on less than 15MB of RAM.
1. Create the Bot Script
On your local machine or a VM, create a file named main.go:
Phase 4: Deployment Options
You can deploy this Golang binary either on a dedicated "Tiny VM" or locally on your iPad using iSH.
Option A: Local Deployment on iPad via iSH (Recommended for portability)
Because iSH runs a 32-bit x86 emulator, we must cross-compile the Go script specifically for that architecture.
1. Compile for iSH (Run this on your main computer):
2. Transfer to iSH:
• Transfer the telegram-bot file to your iPad (e.g., via AirDrop or iCloud Drive).
• Open the iOS Files app, locate the binary, and copy it into the iSH folder under "Locations".
3. Run the Bot in iSH:
• Open the iSH app on your iPad.
• Ensure you have the Tailscale iOS app running and connected on the iPad (iSH uses the iPad's network stack).
• Make the file executable and run it with your environment variables:
Option B: Cloud Deployment on Tiny VM 1
If you want the bot to run 24/7 without keeping iSH open on your iPad, deploy it to a secondary tiny VM.
1. Compile and Transfer:
Move the telegram-bot binary to VM 1.
2. Create the Service File:
Run sudo nano /etc/systemd/system/telegram-bot.service and paste:
3. Enable and Start:
The Result
Your architecture is now complete.
• Open the Telegram app on your iOS 12 iPad.
• Send a message to your bot.
• The message traverses Telegram's servers -> Your Golang binary (iSH or VM 1) -> Tailscale Encrypted Mesh -> LiteLLM Gateway (VM 2) -> Gemini API.
You have successfully bypassed Apple's software obsolescence and brought state-of-the-art AI to a decade-old device in a fully secure, private, and hyper-efficient environment.