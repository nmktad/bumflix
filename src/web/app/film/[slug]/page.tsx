import "@vidstack/react/player/styles/base.css";
import "@vidstack/react/player/styles/plyr/theme.css";

import { MediaPlayer, MediaProvider } from "@vidstack/react";
import {
  PlyrLayout,
  plyrLayoutIcons,
} from "@vidstack/react/player/layouts/plyr";

export default async function FilmViewPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;

  return (
    <MediaPlayer playsInline src={`http://localhost:8080/film/${slug}`}>
      <MediaProvider />
      <PlyrLayout
        thumbnails="https://files.vidstack.io/sprite-fight/thumbnails.vtt"
        icons={plyrLayoutIcons}
      />
    </MediaPlayer>
  );
}
