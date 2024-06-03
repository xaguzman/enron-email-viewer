const apiUrl = import.meta.env.VITE_API_URL;
const apiUser = import.meta.env.VITE_API_USER;
const apiUserPwd = import.meta.env.VITE_API_USER_PWD;

console.log( {apiUrl, apiUser, apiUserPwd } );

export const queryEmails = async (searchQuery: string) => {
    
    const response = await fetch(`${apiUrl}/search`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            query: searchQuery
        })
    });

    let result = [];
    if (response.ok) {
        const data = await response.json();
        result = data
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