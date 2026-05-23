const PLAYING = "▶ Now Playing";
const NOT_PLAYING = "⏸ Not Playing";
const META_SEP = " • ";
const ARTIST_SEP = ", ";
const EVENT_URL = "/events";

let eventSource = new EventSource(EVENT_URL);

const elements = {
  status: document.querySelector("#status"),
  artist: document.querySelector("#artist"),
  title: document.querySelector("#title"),
  release: document.querySelector("#release"),
  meta: document.querySelector("#meta"),
  container: document.querySelector("#player"),
  trackInfo: document.querySelector("#track-info"),
};

const renderNotPlaying = () => {
  elements.status.textContent = NOT_PLAYING;

  elements.title.textContent = "";
  elements.artist.textContent = "";
  elements.release.textContent = "";
  elements.meta.textContent = "";

  elements.container.classList.remove("playing");
  elements.trackInfo.classList.add("hidden");
};

const renderPlaying = (data) => {
  const title = () => {
    if (data.title.length > 0) {
      return data.title;
    }

    return "Unknown Track";
  };

  const artist = () => {
    if (data.artists.length > 0) {
      return data.artists.map((a) => a.name).join(ARTIST_SEP);
    }

    if (data.artist.length > 0) {
      return data.artist;
    }

    return "Unknown Artist";
  };

  const release = () => {
    if (data.release.length > 0) {
      return data.release;
    }

    return "Unknown Release";
  };

  const meta = () => {
    const metaParts = [];

    if (data.trackNumber) {
      metaParts.push(`Track ${data.trackNumber}`);
    }

    if (data.duration) {
      const minutes = Math.floor(data.duration / 60000);
      const seconds = Math.round((data.duration % 60000) / 1000)
        .toString()
        .padStart(2, "0");

      metaParts.push(`${minutes}:${seconds}`);
    }

    return metaParts.join(META_SEP);
  };

  elements.status.textContent = PLAYING;

  elements.title.textContent = title();
  elements.artist.textContent = `by ${artist()}`;
  elements.release.textContent = `on ${release()}`;
  elements.meta.textContent = meta();

  elements.container.classList.add("playing");
  elements.trackInfo.classList.remove("hidden");
};

eventSource.addEventListener("message", (event) => {
  try {
    const data = JSON.parse(event.data);

    console.log("Payload:", data);

    if (data.playing) {
      renderPlaying(data);
    } else {
      renderNotPlaying();
    }
  } catch (err) {
    console.error("Failed to parse payload:", err);
  }
});

eventSource.addEventListener("open", () => {
  console.log("Connection opened");
});

eventSource.addEventListener("error", (event) => {
  console.error("Connection error:", event);

  renderNotPlaying();
});

renderNotPlaying();
