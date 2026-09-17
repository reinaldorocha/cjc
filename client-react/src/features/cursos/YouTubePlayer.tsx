import { useEffect, useRef, useState } from 'react';

type Player = {
  destroy(): void; playVideo(): void; pauseVideo(): void; getCurrentTime(): number;
  getDuration(): number; seekTo(time: number, allowSeekAhead: boolean): void; setVolume(volume: number): void;
};
type YouTubeAPI = { Player: new (element: HTMLElement, options: Record<string, unknown>) => Player };
let apiPromise: Promise<YouTubeAPI> | undefined;
function carregarAPI(): Promise<YouTubeAPI> {
  const ytWindow = window as Window & { YT?: YouTubeAPI; onYouTubeIframeAPIReady?: () => void };
  if (ytWindow.YT?.Player) return Promise.resolve(ytWindow.YT);
  if (apiPromise) return apiPromise;
  apiPromise = new Promise<YouTubeAPI>((resolve, reject) => {
    const timeout = window.setTimeout(() => reject(new Error('O YouTube demorou para responder. Reabra a aula.')), 20000);
    const previous = ytWindow.onYouTubeIframeAPIReady;
    ytWindow.onYouTubeIframeAPIReady = () => {
      previous?.();
      window.clearTimeout(timeout);
      if (ytWindow.YT) resolve(ytWindow.YT);
    };
    const script = document.createElement('script');
    script.src = 'https://www.youtube.com/iframe_api';
    script.onerror = () => { window.clearTimeout(timeout); script.remove(); reject(new Error('Não foi possível conectar ao YouTube. Reabra a aula.')); };
    document.head.appendChild(script);
  }).catch((error: unknown) => { apiPromise = undefined; throw error; });
  return apiPromise;
}
const tempo = (s: number) => `${Math.floor(s / 60)}:${Math.floor(s % 60).toString().padStart(2, '0')}`;

export function YouTubePlayer({ videoId, titulo, inicio = 0, onTempo, onSalvar }: { videoId: string; titulo: string; inicio?: number; onTempo?: (posicao: number, duracao: number) => void; onSalvar?: (posicao: number, duracao: number) => void }) {
  const callbacks = useRef({ onTempo, onSalvar });
  const inicioRef = useRef(inicio);
  useEffect(() => { callbacks.current = { onTempo, onSalvar }; }, [onTempo, onSalvar]);
  const container = useRef<HTMLDivElement>(null);
  const frame = useRef<HTMLDivElement>(null);
  const player = useRef<Player | null>(null);
  const [pronto, setPronto] = useState(false);
  const [tocando, setTocando] = useState(false);
  const [posicao, setPosicao] = useState(0);
  const [duracao, setDuracao] = useState(0);
  const [volume, setVolume] = useState(80);
  const [erro, setErro] = useState('');
  useEffect(() => {
    let ativo = true;
    let instancia: Player | undefined;
    let iniciado = false;
    let reproduzindo = false;
    let ultimoSalvo = Date.now();
    const salvar = () => { if (iniciado && instancia && instancia.getDuration() > 0) callbacks.current.onSalvar?.(instancia.getCurrentTime() || 0, instancia.getDuration()); };
    const ocultar = () => { if (document.visibilityState === 'hidden') salvar(); };
    document.addEventListener('visibilitychange', ocultar);
    window.addEventListener('pagehide', salvar);
    setPronto(false); setTocando(false); setPosicao(0); setDuracao(0); setErro('');
    carregarAPI().then((YT) => {
      if (!ativo || !frame.current) return;
      const alvo = document.createElement('div');
      frame.current.replaceChildren(alvo);
      instancia = new YT.Player(alvo, {
        host: 'https://www.youtube-nocookie.com', videoId,
        playerVars: { controls: 0, disablekb: 1, playsinline: 1, rel: 0, fs: 0, start: Math.floor(inicioRef.current), origin: window.location.origin },
        events: {
          onReady: () => {
            if (!ativo) return;
            player.current = instancia!; instancia!.setVolume(80); setVolume(80); setPronto(true);
            const iframe = frame.current?.querySelector('iframe');
            iframe?.setAttribute('tabindex', '-1'); iframe?.setAttribute('title', titulo);
          },
          onStateChange: (event: { data: number }) => { if (ativo) { if (event.data === 1) iniciado = true; reproduzindo = event.data === 1; setTocando(reproduzindo); if (event.data === 0 || event.data === 2) salvar(); } },
          onError: () => { if (ativo) { setPronto(false); setErro('Esta aula está indisponível ou o vídeo não permite reprodução incorporada. Avise seu mentor.'); } },
        },
      });
    }).catch((e: Error) => { if (ativo) setErro(e.message); });
    const interval = window.setInterval(() => {
      if (player.current) { const p = player.current.getCurrentTime() || 0; const d = player.current.getDuration() || 0; setPosicao(p); setDuracao(d); if (iniciado && d > 0) callbacks.current.onTempo?.(p, d); if (reproduzindo && Date.now() - ultimoSalvo > 15000) { salvar(); ultimoSalvo = Date.now(); } }
    }, 700);
    return () => { salvar(); ativo = false; window.clearInterval(interval); document.removeEventListener('visibilitychange', ocultar); window.removeEventListener('pagehide', salvar); player.current = null; instancia?.destroy(); };
  }, [videoId, titulo]);
  return <div ref={container} className="curso-player" onContextMenu={(e) => e.preventDefault()}>
    <div className="curso-video" ref={frame} />
    {erro && <p role="alert" className="curso-alerta">{erro}</p>}
    {!pronto && !erro && <p role="status">Preparando sua aula…</p>}
    <div className="curso-player-controls">
      <button type="button" disabled={!pronto} onClick={() => tocando ? player.current?.pauseVideo() : player.current?.playVideo()}>{tocando ? 'Ⅱ Pausar' : '▶ Reproduzir'}</button>
      <span>{tempo(posicao)} / {tempo(duracao)}</span>
      <input aria-label="Posição da aula" type="range" min="0" max={duracao || 1} value={posicao} disabled={!pronto} onChange={(e) => { const v = Number(e.target.value); player.current?.seekTo(v, true); setPosicao(v); }} />
      <input aria-label="Volume" className="curso-volume" type="range" min="0" max="100" value={volume} disabled={!pronto} onChange={(e) => { const v = Number(e.target.value); setVolume(v); player.current?.setVolume(v); }} />
      <button type="button" onClick={() => { const action = document.fullscreenElement ? document.exitFullscreen() : container.current?.requestFullscreen(); action?.catch(() => setErro('Tela cheia indisponível neste navegador.')); }}>Tela cheia</button>
    </div>
  </div>;
}
