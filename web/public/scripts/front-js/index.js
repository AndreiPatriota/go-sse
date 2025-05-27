document.addEventListener('alpine:init', () => {
  Alpine.data('sseApp', () => ({
    latestTime: 'Waiting...',
    init() {
      const evtSource = new EventSource('/sse-stream');

      evtSource.onmessage = (event) => {
        const payload = JSON.parse(event.data);
        console.log('Received SSE:', payload);
        this.latestTime = payload.mensagem;
      };

      evtSource.onerror = (err) => {
        console.error('SSE error:', err);
        evtSource.close();
      };
    },
  }));
});
