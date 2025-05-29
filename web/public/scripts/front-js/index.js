document.addEventListener('alpine:init', () => {
  Alpine.data('sseApp', () => ({
    temperatura: 0,
    velocidadeVento: 0,
    init() {
      const evtSource = new EventSource('/sse-stream');

      evtSource.onmessage = (event) => {
        const payload = JSON.parse(event.data);
        console.log('Received SSE:', payload);
        this.temperatura = payload.temperatura || 'No data';
        this.velocidadeVento = payload.velocidade_vento || 'No data';
      };

      evtSource.onerror = (err) => {
        console.error('SSE error:', err);
        evtSource.close();
      };
    },
  }));
});
