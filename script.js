
function wait(delayInMS) {
    return new Promise((resolve) => setTimeout(resolve, delayInMS));
}

// Создание ссылки для скачивания
function createDownloadLink(blob, t) {
    const downloadLink = document.createElement('a');
    downloadLink.href = URL.createObjectURL(blob);
    downloadLink.download = 'recorded-audio.webm';
    downloadLink.textContent = 'Скачать запись';

    // Добавляем на страницу или заменяем существующую
    const existingLink = document.getElementById('downloadLink');
    if (existingLink) {
        existingLink.replaceWith(downloadLink);
    } else {
        downloadLink.id = 'downloadLink';
        t.appendChild(downloadLink);
    }
}

// Панель статуса
class StatusInfo extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        this.shadow = this.attachShadow({ mode: "open" });
        this.cont = this.shadow.appendChild(document.createElement("div"));
        this.cont.style.width = "100%";
        this.cont.style.height = "30px";
        this.cont.style.border = "1px solid red";
        this.cont.style.color = "black";
        this.cont.style.display = "flex";
        this.cont.style.flexDirection = "column";
        this.cont.style.alignItems = "center";
        this.cont.style.justifyContent = "center";
        this.body = this.cont.appendChild(document.createElement("div"));
    }

    refresh() {
        setTimeout(() => {
            this.body.innerHTML = `<div>${window.APP.connection.getJid()}</div>`;
        }, 3000);
    }

}

// Панель для информации о аудио потоке
class AudioInfo extends HTMLElement {
    constructor() {
        super();
        this.displayName = "";
    }

    connectedCallback() {
        console.error("connectedCallback");
        this.shadow = this.attachShadow({ mode: "open" });
        this.cont = this.shadow.appendChild(document.createElement("div"));
        this.cont.style.width = "100%";
        this.cont.style.height = "30px";
        this.cont.style.border = "1px solid green";
        this.cont.style.color = "black";
        this.cont.style.display = "flex";
        this.cont.style.flexDirection = "column";
        this.cont.style.alignItems = "center";
        this.cont.style.justifyContent = "center";
        this.body = this.cont.appendChild(document.createElement("div"));
        this.root = this.cont.appendChild(document.createElement("div"));
    }

    ondataavailable(event) {
        if (event.data.size > 0) {
            console.error(event.data.size);
            const reader = new FileReader();
            reader.onloadend = () => {
                const base64 = reader.result.split(',')[1];
                window.ssbot_writeSound(JSON.stringify({ myid: window.APP.conference.getMyUserId(), room: window.APP.conference.roomName, userid: this.userId, user: this.displayName, u: this.audioElement.id, d: base64 }));
            };
            reader.readAsDataURL(event.data);
        }
    }

    onstop(event) {
        this.root.innerHTML = "stop";
    }

    syncInfo() {
        const lst = window.APP.conference.listMembers();
        for (var m of lst) {
            for (var t of m._tracks) {
                if (t.type === "audio") {
                    for (var c of t.containers) {
                        if (c.id == this.audioElement.id) {
                            this.displayName = m._displayName;
                            this.userId = m._id;
                            this.body.innerHTML = `<b>${this.displayName}</b>`;
                        }
                    }
                }
            }
        }
    }

    init(audioElement) {
        this.audioElement = audioElement;
        this.syncInfo();
        this.root.innerHTML = "init";
        this.recordedChunks = [];
        // Инициализируем AudioContext при первом клике (из-за autoplay policy)
        if (!this.audioContext && !this.initAudioContext()) {
            alert('Не удалось инициализировать аудио контекст');
            return false;
        }

        // Получаем аудиопоток из audio элемента
        this.audioStream = this.audioElement.captureStream();

        // Создаем MediaRecorder
        this.mediaRecorder = new MediaRecorder(this.audioStream, {
            mimeType: 'audio/webm;codecs=opus'
        });

        this.mediaRecorder.onstart = (e) => {
            this.root.innerHTML = "start";
        }
        this.mediaRecorder.ondataavailable = (e) => this.ondataavailable(e);
        this.mediaRecorder.onstop = (e) => this.onstop(e);
        return true;
    }

    initAudioContext() {
        try {
            // Создаем AudioContext
            this.audioContext = new (window.AudioContext || window.webkitAudioContext)();

            // Создаем источник из audio элемента
            const source = this.audioContext.createMediaElementSource(this.audioElement);

            // Подключаем к выходу (динамикам)
            source.connect(this.audioContext.destination);

            console.error('+AudioContext');
            return true;
        } catch (error) {
            console.error('err AudioContext:', error);
            return false;
        }
    }

    startRecording() {
        this.mediaRecorder.start(10000);
    }

    stopRecording() {
        this.mediaRecorder.stop();
    }
}
customElements.define("ssbot-audio", AudioInfo);
customElements.define("ssbot-info", StatusInfo);

// Функция для отслеживания появления/исчезновения элементов
function observeElements(selector, callback) {

    const infopanel = document.body.appendChild(document.createElement("div"));
    infopanel.style.position = "absolute";
    infopanel.style.top = "20px";
    infopanel.style.left = "20px";
    infopanel.style.minHeight = "50px";
    infopanel.style.minWidth = "50px";
    infopanel.style.maxWidth = "400px";
    infopanel.style.padding = "4px";
    infopanel.style.backgroundColor = "white";
    infopanel.innerHTML = "<h5>BotPanel</h5>";
    infopanel.id = "ssbot_panel";

    const ssbotinfo = infopanel.appendChild(document.createElement("ssbot-info"));
    ssbotinfo.refresh();

    const targetNode = document.body;
    const config = {
        childList: true,
        subtree: true
    };

    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            // Проверяем добавленные узлы
            mutation.addedNodes.forEach((node) => {
                if (node.nodeType === 1) { // Проверяем, что это элемент
                    if (node.matches(selector)) {
                        callback(node, 'appeared');
                    }
                    // Проверяем дочерние элементы добавленного узла
                    const matchingElements = node.querySelectorAll(selector);
                    matchingElements.forEach(el => callback(el, 'appeared'));
                }
            });

            // Проверяем удаленные узлы
            mutation.removedNodes.forEach((node) => {
                if (node.nodeType === 1) {
                    if (node.matches(selector)) {
                        callback(node, 'disappeared');
                    }
                    // Проверяем дочерние элементы удаленного узла
                    const matchingElements = node.querySelectorAll(selector);
                    matchingElements.forEach(el => callback(el, 'disappeared'));
                }
            });
        });
    });

    observer.observe(targetNode, config);
    return observer;
}

// Пример использования
const observer = observeElements('audio', (element, event) => {
    if (event === 'appeared') {
        console.error('+:' + element.id);
        if (element.id.startsWith("remoteAudio_")) {
            handleElementAppeared(element);
        }
    } else if (event === 'disappeared') {
        console.error('-:' + element.id);
        if (element.id.startsWith("remoteAudio_")) {
            handleElementDisappeared(element);
        }
    }
});

let audios = {};

function handleElementAppeared(element) {
    const i = document.getElementById("ssbot_panel").appendChild(document.createElement("ssbot-audio"));
    audios[element.id] = i;
    i.init(element);
    i.startRecording();
}

function handleElementDisappeared(element) {
    console.error("mr stop: " + audios[element.id].mediaRecorder.state);
    audios[element.id].stopRecording();
    audios[element.id].remove();
}
"";