export type User = {
  id: number;
  gitlabEmail: string;
  gitlabUsername: string; // currently not filled with data
  gitlabAvatar: { avatarURL: string; fallbackAvatarURL: string }; // currently not filled with data
  name: string;
};
