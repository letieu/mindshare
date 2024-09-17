document.getElementById("share-button").addEventListener("click", async () => {
  if (navigator.share) {
    try {
      await navigator.share({
        title: document.title,
        url: window.location.href,
      });
      console.log("Successfully shared");
    } catch (error) {
      console.error("Error sharing:", error);
    }
  } else {
    // fallback
    console.log("Web Share API not supported");
    navigator.clipboard.writeText(window.location.href);
    alert("Link copied to clipboard");
  }
});
