document.addEventListener('alpine:init', () => {
  Alpine.data('sseApp', () => ({
    temperatura: 0,
    velocidadeVento: 0,
    umidadeRelativa: 0,
    hora: '00:00',
    init() {
      const evtSource = new EventSource('/sse-stream');

      evtSource.addEventListener('update', (event) => {
        const jsonPayload = JSON.parse(event.data);
        console.log('Received SSE:', jsonPayload);
        this.temperatura = jsonPayload.temperatura || 'No data';
        this.velocidadeVento = jsonPayload.velocidade_vento || 'No data';
        this.umidadeRelativa = jsonPayload.umidade_relativa || 'No data';
        this.hora = jsonPayload.hora || 'No data';
      });

      evtSource.addEventListener('error', (event) => {
        const jsonPayload = JSON.parse(event.data);
        window.alert(`SSE error: ${jsonPayload.mensagem}`);
        evtSource.close();
      });
    },
  }));
});
