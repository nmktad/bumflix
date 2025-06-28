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
    <MediaPlayer
      playsInline
      src={`${process.env.NEXT_PUBLIC_BACKEND_URL}/film/${slug}`}
    >
      <MediaProvider />
      <PlyrLayout icons={plyrLayoutIcons} />
    </MediaPlayer>
  );
}
