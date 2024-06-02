const apiUrl = import.meta.env.VITE_API_URL;
const apiUser = import.meta.env.VITE_API_USER;
const apiUserPwd = import.meta.env.VITE_API_USER_PWD;

console.log( {apiUrl, apiUser, apiUserPwd } );

export const queryEmails = async (searchQuery: string) => {
    
    const response = await fetch(`${apiUrl}/api/enron-emails/_search`, {
        method: 'POST',
        headers: {

            'Content-Type': 'application/json',
            'Authorization': `Basic ${ encodedCredentials() }`
        },
        body: JSON.stringify({
            search_type: "matchphrase",
            query: {
                term: searchQuery,
                field: "_all"
            },
            sort_fields: ["Date"],
            from: 0,
            max_results: 20,
            highlight: {
                fields: {
                    "Body": {
                        "pre_tags": ["<mark class=\"highlight\">"],
                        "post_tags": ["</mark>"]
                    },
                    "Subject": {
                        "pre_tags": ["<mark class=\"highlight\">"],
                        "post_tags": ["</mark>"]
                    }
                }
            }
        })
    });

    let result = [];
    if (response.ok) {
        const data = await response.json();
        result = data.hits.hits
            .map( (hit : any)  => {
                let result = hit._source
                result.Id = hit._id;
                return result;
            });
        
    } else {
        console.error('Failed to fetch emails:', response.statusText);
    }

    return result;
}

const encodedCredentials = () => {
    return btoa(`${apiUser}:${apiUserPwd}`);
}