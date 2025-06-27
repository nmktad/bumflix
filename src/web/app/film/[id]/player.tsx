"use client";

import "@vidstack/react/player/styles/base.css";
import "@vidstack/react/player/styles/plyr/theme.css";

import { MediaPlayer, MediaProvider, Captions } from "@vidstack/react";
import {
  PlyrLayout,
  plyrLayoutIcons,
} from "@vidstack/react/player/layouts/plyr";

import { Film } from "@/types/types";

export default function Player({ film }: { film: Film }) {
  return (
    <MediaPlayer
      playsInline
      src="https://files.vidstack.io/sprite-fight/hls/stream.m3u8"
    >
      <MediaProvider />
      <PlyrLayout
        thumbnails="https://files.vidstack.io/sprite-fight/thumbnails.vtt"
        title={film.title}
        icons={plyrLayoutIcons}
      />
      <Captions className="media-captions absolute inset-0 bottom-2 z-10 select-none break-words opacity-0 transition-[opacity,bottom] duration-300 media-captions:opacity-100 media-controls:bottom-[85px] media-preview:opacity-0 aria-hidden:hidden" />
    </MediaPlayer>
  );
}
