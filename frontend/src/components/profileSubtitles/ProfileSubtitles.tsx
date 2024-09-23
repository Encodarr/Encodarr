import InputCheckbox from "../inputs/inputCheckbox/InputCheckbox";
import InputContainer from "../inputs/inputContainer/InputContainer";
import styles from "./ProfileSubtitles.module.scss";

const ProfileSubtitles = ({ content, setContent }) => {
    const languages = {
        all: "All",
        eng: "English",
        spa: "Spanish",
        fra: "French",
        deu: "German",
        ita: "Italian",
        jpn: "Japanese",
        kor: "Korean",
        rus: "Russian",
        zho: "Chinese",
    };
    return (
        <>
            <div className={styles.section}>
                <div className={styles.left}>
                    <InputContainer
                        label="Include Untagged Subtitle tracks"
                        type="checkbox"
                        checked={content?.mapUntaggedSubtitleTracks}
                        onChange={(e) => {
                            setContent({
                                ...content,
                                mapUntaggedSubtitleTracks: e.target.checked,
                            });
                        }}
                    />
                </div>
                <div className={styles.right}>
                    <label>Wanted Subtitle Tracks</label>
                    <div className={styles.targets}>
                        {Object.entries(languages)?.map(([key, value]: any) => (
                            <div
                                key={key}
                                className={styles.target}
                                style={
                                    !content?.profileSubtitleLanguages?.some(lang => lang.language === key)
                                        ? { opacity: "50%" }
                                        : {}
                                }
                                onClick={() => {
                                    const selectedValues = [...content.profileSubtitleLanguages];
                                    const index = selectedValues.findIndex(lang => lang.language === key);
                                    if (index > -1) {
                                        selectedValues.splice(index, 1);
                                    } else {
                                        selectedValues.push({ profileId: content.id, language: key });
                                    }
                                    setContent({ ...content, profileSubtitleLanguages: selectedValues });
                                }}
                            >
                                <InputCheckbox
                                    type="checkbox"
                                    checked={content?.profileSubtitleLanguages?.some(lang => lang.language === key)}
                                    onChange={() => {
                                        const selectedValues = [...content.profileSubtitleLanguages];
                                        const index = selectedValues.findIndex(lang => lang.language === key);
                                        if (index > -1) {
                                            selectedValues.splice(index, 1);
                                        } else {
                                            selectedValues.push({ profileId: content.id, language: key });
                                        }
                                        setContent({ ...content, profileSubtitleLanguages: selectedValues });
                                    }}
                                />
                                <span className={styles.key}>{value}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </>
    );
};

export default ProfileSubtitles;