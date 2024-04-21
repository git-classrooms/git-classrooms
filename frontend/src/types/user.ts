export type User = {
  id: number;
  gitlabEmail: string;
  gitlabUrl: string; // currently not filled with data
  gitlabUsername: string; // currently not filled with data
  gitlabAvatar: { avatarURL: string; fallbackAvatarURL: string }; // currently not filled with data
  name: string;
};
