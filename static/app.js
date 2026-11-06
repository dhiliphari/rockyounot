document.addEventListener('alpine:init', () => {
  Alpine.data('passwordForm', (endpoint) => ({
    password: '',
    show: false,
    loading: false,
    result: null,
    error: '',

    async submit() {
      this.loading = true;
      this.result = null;
      this.error = '';

      try {
        const res = await fetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ password: this.password }),
        });

        if (!res.ok) {
          this.error = (await res.text()).trim() || `Request failed (${res.status})`;
          return;
        }

        const body = await res.text();
        this.result = body ? JSON.parse(body) : {};
      } catch (err) {
        this.error = 'Could not reach the server.';
      } finally {
        this.loading = false;
      }
    },
  }));
});
