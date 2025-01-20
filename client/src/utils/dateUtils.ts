export const formatDateTime = (dateString: string, locale?: string): string => {
  const date = new Date(dateString);
  const userLocale = locale || navigator.language || "en-US";

  return new Intl.DateTimeFormat(userLocale, {
    dateStyle: "medium",
    timeStyle: "medium",
  }).format(date);
};
