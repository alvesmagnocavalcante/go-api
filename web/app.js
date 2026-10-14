const form = document.querySelector('#search-form');
const input = document.querySelector('#name');
const submit = document.querySelector('#submit');
const status = document.querySelector('#status');
const result = document.querySelector('#result');
const suggestions = document.querySelectorAll('[data-name]');
const labels = { hp: 'HP', attack: 'Ataque', defense: 'Defesa', 'special-attack': 'Ataque especial', 'special-defense': 'Defesa especial', speed: 'Velocidade' };
const types = { normal: 'Normal', fire: 'Fogo', water: 'Água', electric: 'Elétrico', grass: 'Planta', ice: 'Gelo', fighting: 'Lutador', poison: 'Veneno', ground: 'Terra', flying: 'Voador', psychic: 'Psíquico', bug: 'Inseto', rock: 'Pedra', ghost: 'Fantasma', dragon: 'Dragão', dark: 'Sombrio', steel: 'Aço', fairy: 'Fada' };
const decimal = new Intl.NumberFormat('pt-BR', { maximumFractionDigits: 1 });

function render(data) {
  document.querySelector('#number').textContent = `#${String(data.id).padStart(3, '0')}`;
  document.querySelector('#pokemon-name').textContent = data.name;
  document.querySelector('#height').textContent = `${decimal.format(data.height / 10)} m`;
  document.querySelector('#weight').textContent = `${decimal.format(data.weight / 10)} kg`;
  const image = document.querySelector('#sprite');
  const missing = document.querySelector('#no-image');
  const url = data.sprites?.other?.['official-artwork']?.front_default || data.sprites?.front_default;
  image.hidden = !url;
  missing.hidden = Boolean(url);
  image.onerror = () => { image.hidden = true; missing.hidden = false; };
  image.alt = data.name;
  if (url) image.src = url;
  else image.removeAttribute('src');
  document.querySelector('#types').replaceChildren(...data.types.map(({ type }) => {
    const badge = document.createElement('span');
    badge.className = 'type';
    badge.textContent = types[type.name] || type.name;
    return badge;
  }));
  document.querySelector('#stats').replaceChildren(...data.stats.map(({ stat, base_stat }) => {
    const row = document.createElement('div');
    row.className = 'stat';
    const label = document.createElement('span');
    label.textContent = labels[stat.name] || stat.name;
    const meter = document.createElement('meter');
    meter.min = 0;
    meter.max = 255;
    meter.value = base_stat;
    meter.setAttribute('aria-label', label.textContent);
    const value = document.createElement('strong');
    value.textContent = base_stat;
    row.append(label, meter, value);
    return row;
  }));
  result.hidden = false;
}

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  const name = input.value.trim().toLowerCase();
  if (!name) return;
  submit.disabled = true;
  suggestions.forEach(button => { button.disabled = true; });
  result.hidden = true;
  form.setAttribute('aria-busy', 'true');
  status.classList.remove('error');
  status.textContent = 'Buscando Pokémon…';
  try {
    const response = await fetch(`/api/pokemon/${encodeURIComponent(name)}`, { signal: AbortSignal.timeout(12000) });
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || 'Não foi possível realizar a busca.');
    render(data);
    status.textContent = `${data.name} encontrado!`;
  } catch (error) {
    result.hidden = true;
    status.classList.add('error');
    status.textContent = error.name === 'TimeoutError' ? 'A consulta demorou demais. Tente novamente.' : error instanceof TypeError ? 'Falha na conexão. Tente novamente.' : error.message;
  } finally {
    submit.disabled = false;
    suggestions.forEach(button => { button.disabled = false; });
    form.setAttribute('aria-busy', 'false');
  }
});

suggestions.forEach(button => button.addEventListener('click', () => {
  input.value = button.dataset.name;
  form.requestSubmit();
}));
